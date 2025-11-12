package cli

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/BehnamAxo/driftq-cli/pkg/driftqclient"
)

const defaultAddr = "localhost:9090"

// This whole thing needs refactoring

func Run(args []string) error {
	if len(args) == 0 {
		printHelp()
		return nil
	}

	switch args[0] {
	case "help", "-h", "--help":
		printHelp()
		return nil
	case "health":
		return cmdHealth(args[1:])
	case "topic":
		return cmdTopic(args[1:])
	case "produce":
		return cmdProduce(args[1:])
	case "consume":
		return cmdConsume(args[1:])
	case "lag":
		return cmdLag(args[1:])
	default:
		return fmt.Errorf("unknown command %q; run 'driftq help'", args[0])
	}
}

// TODO: Test this Behnam
func printHelp() {
	fmt.Println(`driftq CLI
Usage:
  driftq health [--addr HOST:PORT]
  driftq topic create <name> [--partitions N] [--compacted] [--addr HOST:PORT]
  driftq topic list [--addr HOST:PORT]
  driftq produce <topic> [--key K] [--value V|--file PATH] [--addr HOST:PORT]
  driftq consume <topic> --group G [--from latest|earliest] [--addr HOST:PORT]
  driftq lag <topic> [--group G] [--addr HOST:PORT]

Notes:
  - This CLI talks to the broker at --addr. Current build uses stub calls until the admin/produce APIs are live.
`)
}

func newClient(addr string) *driftqclient.Client {
	if addr == "" {
		addr = defaultAddr
	}
	return driftqclient.New(addr)
}

func cmdHealth(args []string) error {
	fs := flag.NewFlagSet("health", flag.ExitOnError)
	addr := fs.String("addr", defaultAddr, "broker address host:port")
	_ = fs.Parse(args)
	c := newClient(*addr)
	ok, err := c.Health()
	if err != nil {
		return err
	}
	if ok {
		fmt.Println("OK")
	} else {
		fmt.Println("NOT OK")
	}
	return nil
}

func cmdTopic(args []string) error {
	if len(args) == 0 {
		return errors.New("usage: driftq topic [create|list] ...")
	}
	switch args[0] {
	case "create":
		fs := flag.NewFlagSet("topic create", flag.ExitOnError)
		addr := fs.String("addr", defaultAddr, "broker addr")
		parts := fs.Int("partitions", 1, "number of partitions")
		compacted := fs.Bool("compacted", false, "enable log compaction")
		fs.Parse(args[1:])
		if fs.NArg() < 1 {
			return errors.New("missing <name>")
		}
		name := fs.Arg(0)
		c := newClient(*addr)
		return c.TopicCreate(name, *parts, *compacted)
	case "list":
		fs := flag.NewFlagSet("topic list", flag.ExitOnError)
		addr := fs.String("addr", defaultAddr, "broker addr")
		fs.Parse(args[1:])
		c := newClient(*addr)
		topics, err := c.TopicList()
		if err != nil {
			return err
		}
		for _, t := range topics {
			fmt.Printf("%s\tpartitions=%d\tcompacted=%v\n", t.Name, t.Partitions, t.Compacted)
		}
		return nil
	default:
		return fmt.Errorf("unknown topic subcommand %q", args[0])
	}
}

func cmdProduce(args []string) error {
	fs := flag.NewFlagSet("produce", flag.ExitOnError)
	addr := fs.String("addr", defaultAddr, "broker addr")
	key := fs.String("key", "", "message key")
	value := fs.String("value", "", "message value (string)")
	file := fs.String("file", "", "read value from file path")
	fs.Parse(args)
	if fs.NArg() < 1 {
		return errors.New("missing <topic>")
	}
	topic := fs.Arg(0)
	if *value != "" && *file != "" {
		return errors.New("use either --value or --file, not both")
	}
	var data []byte
	if *file != "" {
		b, err := os.ReadFile(*file)
		if err != nil {
			return err
		}
		data = b
	} else {
		data = []byte(*value)
	}
	c := newClient(*addr)
	return c.Produce(topic, *key, data)
}

func cmdConsume(args []string) error {
	fs := flag.NewFlagSet("consume", flag.ExitOnError)
	addr := fs.String("addr", defaultAddr, "broker addr")
	group := fs.String("group", "", "consumer group (required)")
	from := fs.String("from", "latest", "start offset: latest|earliest")
	fs.Parse(args)
	if fs.NArg() < 1 {
		return errors.New("missing <topic>")
	}
	if *group == "" {
		return errors.New("--group is required")
	}
	topic := fs.Arg(0)
	c := newClient(*addr)
	return c.Consume(topic, *group, *from, func(key string, value []byte) error {
		fmt.Printf("%s\t%s\n", key, strings.TrimSpace(string(value)))
		return nil
	})
}

func cmdLag(args []string) error {
	fs := flag.NewFlagSet("lag", flag.ExitOnError)
	addr := fs.String("addr", defaultAddr, "broker addr")
	group := fs.String("group", "", "consumer group (optional)")
	fs.Parse(args)
	if fs.NArg() < 1 {
		return errors.New("missing <topic>")
	}
	topic := fs.Arg(0)
	c := newClient(*addr)
	lag, err := c.Lag(topic, *group)
	if err != nil {
		return err
	}
	fmt.Printf("lag=%d (stub)\n", lag)
	return nil
}

func parseDuration(s string) (time.Duration, error) { return time.ParseDuration(s) }
