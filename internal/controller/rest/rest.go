package rest

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/RenterRus/dwld-bot/internal/repo/dwld"
	"github.com/RenterRus/dwld-bot/internal/repo/persistent"
)

type ServerConf struct {
	DB   persistent.SQLRepo
	Host string
	Port int
}

type Server struct {
	db      persistent.SQLRepo
	servers []dwld.DWLDModel
	host    string
	port    int
}

func NewServer(conf *ServerConf) *Server {
	return &Server{
		db:   conf.DB,
		host: conf.Host,
		port: conf.Port,
	}
}

func (s *Server) getServers() ([]dwld.DWLDModel, error) {
	servers, err := s.db.LoadServers()
	if err != nil {
		return nil, fmt.Errorf("servers: %w", err)
	}

	resp := make([]dwld.DWLDModel, 0, len(servers))
	for _, v := range servers {
		resp = append(resp, dwld.NewDWLD(v.Host, v.Port).SetName(v.Name))
	}

	return resp, nil
}

type Links struct {
	Link []string `json:"links"`
}

func (s *Server) list(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Printf("[ERROR] ReadAll: %s\n", err.Error())

		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad request"))
	}

	defer r.Body.Close()

	links := Links{}

	err = json.Unmarshal(body, &links)
	if err != nil {
		fmt.Printf("[ERROR] Unmarshal: %s\n", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
	}

	s.servers, err = s.getServers()
	if err != nil {
		fmt.Printf("[ERROR] getServers: %s\n", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal error"))
	}

	respRaw := Links{
		Link: make([]string, 0, len(links.Link)),
	}

	if len(s.servers) > 0 && len(links.Link) > 0 {
		fmt.Println("servers:")
		for _, srv := range s.servers {
			fmt.Println(srv.Name())
		}
		fmt.Println("---------")

		for _, server := range s.servers {
			for _, l := range links.Link {
				if _, err = url.ParseRequestURI(l); err != nil {
					continue
				}

				server.SetToQueue(context.Background(), l, "sandbox", 10000)

				fmt.Printf("send to download: %s\n", l)

				respRaw.Link = append(respRaw.Link, l)
			}
		}
	}

	resp, err := json.Marshal(respRaw)
	if err != nil {
		fmt.Printf("[ERROR] Marshal: %s\n", err.Error())

		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Queue sended but marshal response failed"))
	}

	w.WriteHeader(http.StatusOK)
	w.Write(resp)
}

func (s *Server) Process() {
	mux := http.NewServeMux()
	mux.HandleFunc("/link", s.list)

	fmt.Printf("listen: %s:%d/link\n", s.host, s.port)

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", s.host, s.port), mux); err != nil {
		log.Fatal(err)
	}
}
