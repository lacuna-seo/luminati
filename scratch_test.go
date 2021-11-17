package luminati

import (
	"fmt"
	"log"
)

func (t *LuminatiTestSuite) Test_All() {
	t.T().Skip()
	client, err := New()
	if err != nil {
		log.Fatalln(err)
	}

	serps, meta, err := client.JSON(Options{
		Keyword: "macbook",
		Country: "us",
		Params:  nil,
		Desktop: false,
	})
	if err != nil {
		log.Fatalln(err)
	}

	domain := serps.CheckURL("https://www.apple.com")

	fmt.Printf("Domain: %+v\n", domain)
	fmt.Printf("Meta: %+v\n", meta)
}
