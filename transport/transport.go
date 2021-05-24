package transport

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"github.com/QuanAnh/rabbitmq-go/endpoints"
	amqptransport "github.com/go-kit/kit/transport/amqp"
	"github.com/streadway/amqp"
	"log"
	"sync"
)

func MakeAMQPHandle(siteEndpoints endpoints.Endpoints, url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	defer conn.Close()
	ch, err := conn.Channel()
	if err != nil {
		return err
	}
	defer ch.Close()

	options := []amqptransport.SubscriberOption{
		amqptransport.SubscriberBefore(),
		amqptransport.SubscriberErrorEncoder(encodeError),
	}

	userQueue, err := ch.QueueDeclare(
		"user-queue",
		false,
		false,
		false,
		false,
		nil,
		)

	// Declare binding key
	err = ch.QueueBind(
		userQueue.Name, // queue name
		"user.get-user-by-id",    // routing key
		"amq.topic",  // exchange
		false,
		nil,
	)
	if err != nil {
		log.Println("lỗi")
	}

	replyQueue, err := ch.QueueDeclare(
		"user-publisher",
		false, // durable
		false, // autoDelete
		false, // exclusive
		false, // noWait
		nil,   //args
	)

	userMsgs, err := ch.Consume(
		userQueue.Name,
		"",    // consumer
		false, // autoAck
		false, // exclusive
		false, // noLocal
		false, // noWait
		nil,   // args
	)

	getUserByIdHandler := amqptransport.NewSubscriber(
		siteEndpoints.GetUserById,
		decodeUserRequest,
		encodeResponse,
		options...,
	)

	getUserByIdListener := getUserByIdHandler.ServeDelivery(ch)

	getUserByIdPublisher := amqptransport.NewPublisher(
		ch,
		&replyQueue,
		encodeGetUserByIdAMQPRequest,
		decodeGetUserByIdAMQPResponse,
		amqptransport.PublisherBefore(
			// queue name specified by subscriber
			amqptransport.SetPublishKey(userQueue.Name),
		),
	)
	getUserByIdEndpoint := getUserByIdPublisher.Endpoint()

	/*go func() {
		for msg := range userMsgs {
			fmt.Sprintf("Received message: %s", msg.RoutingKey)

			var nack bool

			// XXX: Instrumentation to be added in a future episode

			// XXX: We will revisit defining these topics in a better way in future episodes
			switch msg.RoutingKey {
			case "user.get-user-by-id" :
				userRequest := endpoints.UserRequest{Id: "123"}
				response, err := getUserByIdEndpoint(context.Background(), userRequest)
				if err != nil {
					log.Println("error with user request", err)
				} else {
					log.Println("user response: ", response)
				}
			default:
				nack = true
			}
			if nack {
				msg.Nack(false, nack)
			} else {
				msg.Ack(false)
			}
		}
	}()*/


	var wg sync.WaitGroup

	word := flag.String(
		"Id",
		"abc",
		"Word to input",
	)
	userRequest := endpoints.UserRequest{Id: *word}
	wg.Add(1)
	go func() {
		log.Println("sending user request")
		defer wg.Done()
		response, err := getUserByIdEndpoint(context.Background(), userRequest)
		if err != nil {
			log.Println("error with user request", err)
		} else {
			log.Println("user response: ", response)
		}
	}()

	forever := make(chan bool)
	go func() {
		for true {
			select {
			case getUserByIdDeliv := <-userMsgs:
				log.Println("received user request")
				getUserByIdListener(&getUserByIdDeliv)
				getUserByIdDeliv.Ack(false) // multiple = false
			}
		}
	}()

	log.Println("listening")
	<-forever
	wg.Wait()

	return nil
}

var badRequest = errors.New("bad request type")

func encodeError(ctx context.Context, err error, deliv *amqp.Delivery, ch amqptransport.Channel, pub *amqp.Publishing) {

}

func decodeUserRequest(ctx context.Context, d *amqp.Delivery) (interface{}, error) {
	req := endpoints.UserRequest{}
	err := json.Unmarshal(d.Body, &req)
	return req, err
}

func encodeResponse(_ context.Context, pub *amqp.Publishing, response interface{}) error {
	data, _ := json.Marshal(endpoints.UserResponse{
		V:   "ok",
		Err: "",
	})
	pub.Body = data
	return nil
}

func encodeGetUserByIdAMQPRequest(ctx context.Context, publishing *amqp.Publishing, request interface{}) error {
	userRequest, ok := request.(endpoints.UserRequest)
	if !ok {
		return badRequest
	}
	b, err := json.Marshal(userRequest)
	if err != nil {
		return err
	}
	publishing.Body = b
	return nil
}

func decodeGetUserByIdAMQPResponse(ctx context.Context, delivery *amqp.Delivery) (interface{}, error) {
	var response endpoints.UserResponse
	err := json.Unmarshal(delivery.Body, &response)
	if err != nil {
		return nil, err
	}
	return response, nil
}