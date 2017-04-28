package main

import (
	"log"
	"math/rand"
	"strconv"
	"time"

	"github.com/Shopify/sarama"
	"github.com/jessevdk/go-flags"
)

type CommandOptions struct {
	Name []*string `long:"name" description:"name of the sensor" required:"true"`
	Freq int       `short:"f" long:"freq" description:"update frequency in cycles/sec" required:"true"`
	Max  float64   `long:"max" description:"maximum value for generated readings" required:"true"`
	Min  float64   `long:"min" description:"minimum value for generated readings" required:"true"`
	Step float64   `short:"s" long:"step" description:"maximum allowable change per measurement" required:"true"`
	Dir  string    `short:"d" long:"dir" description:"direction of simulation" required:"true"` // can be up, down or other
}

// TemperatureData is the type we send to kafka
type TemperatureData struct {
	SensorName  string
	Temperature float64
	TimeStamp   time.Time
}

type quickSerializer interface {
	toString() string
}

func (temperatureData *TemperatureData) toString() string {
	temperature := strconv.FormatFloat(temperatureData.Temperature, 'f', 2, 64)
	return "{\"sensorName\" : \"" + temperatureData.SensorName + "\"," +
		"\"temperature\" : \"" + temperature + "\"," +
		"\"timeStamp\" : \"" + temperatureData.TimeStamp.Format("2006-01-02 15:04:05.999999999 -0700 MST") + "\"}"
}

func main() {
	start := time.Now()

	options := CommandOptions{}

	parser := initParser(&options)

	_, err := parser.Parse()

	if err != nil {
		panic(err)
	}

	run(&options)

	log.Print(time.Since(start))
}

func init() {
	var banner = `
  |------------------------------------------------------------------------------|
  |                                                                              |
  |              SensorStreaming Golang - Create and Publish sensor data         |
  |                                                                              |
  |------------------------------------------------------------------------------|
  `

	log.Print(banner)
}

func initParser(options *CommandOptions) (parser *flags.Parser) {
	//default behaviour is HelpFlag | PrintErrors | PassDoubleDash - we need to override the stderr output
	return flags.NewParser(options, flags.HelpFlag)
}

func run(options *CommandOptions) {
	log.Print(options)

	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	var value = r.Float64()*(options.Max-options.Min) + options.Min
	var nom = (options.Max-options.Min)/2 + options.Min

	duration, _ := time.ParseDuration(strconv.Itoa(1000/int(options.Freq)) + "ms")
	signal := time.Tick(duration)

	i := 1.
	for range signal {
		produce("platformTemperatureTopic", generateValue(&i, value, nom, r, options))
	}
}

func generateValue(i *float64, value, nom float64, r *rand.Rand, options *CommandOptions) []TemperatureData {
	temperatureDatas := make([]TemperatureData, 0)

	for _, sensorName := range options.Name {
		if options.Dir == "other" {
			temperatureDatas = append(temperatureDatas, generateRandomValue(value, nom, *sensorName, r, options))
		} else if options.Dir == "up" {
			*i += 0.2
			temperatureDatas = append(temperatureDatas, generateFunctionValue(*i, nom, *sensorName, r, options))
		} else if options.Dir == "down" {
			*i -= 0.2
			temperatureDatas = append(temperatureDatas, generateFunctionValue(*i, nom, *sensorName, r, options))
		}
	}
	return temperatureDatas
}

func generateRandomValue(value, nom float64, sensorName string, r *rand.Rand, options *CommandOptions) TemperatureData {
	var maxStep, minStep float64

	if value < nom {
		maxStep = options.Step
		minStep = -1 * options.Step * (value - options.Min) / (nom - options.Min)
	} else {
		maxStep = options.Step * (options.Max - value) / (options.Max - nom)
		minStep = -1 * options.Step
	}

	value += r.Float64()*(maxStep) + r.Float64()*(maxStep-minStep) + minStep
	timeStamp := time.Now()

	temperatureData := TemperatureData{}
	temperatureData.SensorName = sensorName
	temperatureData.Temperature = value
	temperatureData.TimeStamp = timeStamp

	return temperatureData
}

func generateFunctionValue(i, nom float64, sensorName string, r *rand.Rand, options *CommandOptions) TemperatureData {
	tempVal := i + ((r.Float64()) * (options.Max - options.Min))
	nom += tempVal

	temperatureData := TemperatureData{}
	temperatureData.SensorName = sensorName
	temperatureData.Temperature = nom
	temperatureData.TimeStamp = time.Now()

	return temperatureData
}

func produce(topic string, temperatures []TemperatureData) {

	config := sarama.NewConfig()
	brokers := []string{"localhost:9092"}

	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		// Should not reach here
		panic(err)
	}

	defer func() {
		if err := producer.Close(); err != nil {
			// Should not reach here
			panic(err)
		}
	}()

	for i := 0; i < len(temperatures)-1; i++ {

		strTime := strconv.Itoa(int(time.Now().Unix()))
		msg := &sarama.ProducerMessage{
			Topic: topic,
			Key:   sarama.StringEncoder(strTime),
			Value: sarama.StringEncoder(sarama.StringEncoder(temperatures[i].toString())),
		}

		partition, offset, err := producer.SendMessage(msg)
		if err != nil {
			panic(err)
		}

		log.Printf("Message is stored in topic(%s)/partition(%d)/offset(%d)\n", topic, partition, offset)
	}
}
