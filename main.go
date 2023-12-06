package main

import (
	"context"
	"fmt"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {
	cfg, err := Get()
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get config")
	}

	db, err := NewDatabaseInstance(cfg)
	if err != nil {
		log.Fatal().Err(err).Msg("failed to get database instance")
	}

	fmt.Println(db)

	ctx := context.Background()

	// create a slice of names
	// names := []string{"Alice", "Bob", "Charlie", "David", "Eve", "Frank", "Grace", "Harry"}

	dcRepo := NewDonutRepository(db)
	dc := NewDonutCall(dcRepo)

	matchMaker := &MatchMakerEntity{}
	matchMaker.Build(
		WithMatchMakerEntityName("test"),
		WithMatchMakerEntityDescription("test description"),
		WithMatchMakerEntityStartTime(time.Now()),
		WithMatchMakerEntityDuration(10*24*time.Hour),
	)

	dc.CreateMatchMaker(ctx, matchMaker)

	<-time.After(10 * time.Second)

	// for _, name := range names {
	// 	dc.Register(name)
	// }

	// fmt.Println()
	// fmt.Println("Start...")

	// dc.Start()

	// fmt.Println()
	// fmt.Println("Do calls...")

	// dc.DoCall(dc.GetPerson("Alice"), dc.GetPerson("Bob"))
	// dc.DoCall(dc.GetPerson("Charlie"), dc.GetPerson("David"))
	// dc.DoCall(dc.GetPerson("Eve"), dc.GetPerson("Frank"))

	// fmt.Println()
	// fmt.Println("Add person...")

	// dc.AddPerson("Ivan")
	// dc.AddPerson("Goldi")
	// dc.AddPerson("Samde")

	// fmt.Println()
	// fmt.Println("RePair...")

	// dc.RePair()

	// fmt.Println()
	// fmt.Println("Remove person...")

	// dc.RemovePerson("Ivan")

	// fmt.Println()
	// fmt.Println("RePair...")

	// dc.RePair()

	// fmt.Println()
	// fmt.Println("Print...")

	// dc.Print()
}
