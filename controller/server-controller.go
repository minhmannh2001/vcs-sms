package controller

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os/exec"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/go-mail/mail"
	"github.com/minhmannh2001/sms/cache"
	"github.com/minhmannh2001/sms/entity"
	"github.com/minhmannh2001/sms/service"
	elastic "github.com/olivere/elastic/v7"
	"github.com/xuri/excelize/v2"
	"gopkg.in/matryer/try.v1"
)

type time_event_element struct {
	from     time.Time
	to       time.Time
	duration time.Duration
}

type ServerController interface {
	CreateServer(server *entity.Server) error
	// CreateServers()
	ViewServer(id int) (entity.Server, error)
	ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error)
	UpdateServer(server *entity.Server) error
	DeleteServer(id int) error
	ImportServers()
	ExportServers()
	ReportServerStatus()
	AutoReportServerStatus()
	CheckServerExistence(ip string) bool
	CheckServerName(name string) bool
}

type serverController struct {
	serverService service.ServerService
	serverCache   cache.ServerCache
}

func NewServerController(service service.ServerService, cache cache.ServerCache) ServerController {
	return &serverController{
		serverService: service,
		serverCache:   cache,
	}
}

func (controller *serverController) CreateServer(server *entity.Server) error {
	out, _ := exec.Command("ping", server.Ipv4, "-c 5", "-w 10").Output()
	if strings.Contains(string(out), "100% packet loss") {
		fmt.Println("Connot connect to server having ip address " + server.Ipv4)
		return errors.New("Cannot connect to server")
	} else {
		server.Status = "Up"
		fmt.Println("Successfully ping to server having ip address " + server.Ipv4)
		err := controller.serverService.CreateServer(server)
		if err != nil {
			return err
		}
		return nil
	}
}

func (controller *serverController) ViewServer(id int) (entity.Server, error) {
	var server *entity.Server = controller.serverCache.Get(fmt.Sprint(id))
	if server == nil {
		server, err := controller.serverService.ViewServer(id)
		controller.serverCache.Set(fmt.Sprint(id), &server)
		return server, err
	}
	return *server, nil
}

func (controller *serverController) ViewServers(from int, to int, perpage int, sortby string, order string, filter string) ([]entity.Server, error) {
	log.Println("hello")
	key := strconv.FormatInt(int64(from+2), 10) + "-" + strconv.FormatInt(int64(to+2), 10) + "-" + strconv.FormatInt(int64(perpage+2), 10) + "-" + sortby + "-" + order + "-" + filter
	var servers []entity.Server = controller.serverCache.AGet(key)
	if servers == nil {
		result, err := controller.serverService.ViewServers(from, to, perpage, sortby, order, filter)
		controller.serverCache.ASet(key, result)
		return result, err
	}
	return controller.serverService.ViewServers(from, to, perpage, sortby, order, filter)
}

func (controller *serverController) UpdateServer(server *entity.Server) error {
	err := controller.serverService.UpdateServer(server)
	if err != nil {
		return err
	}
	return nil
}

func (controller *serverController) DeleteServer(id int) error {
	err := controller.serverService.DeleteServer(id)
	return err
}

func (controller *serverController) ImportServers() {
	fmt.Println("Execute ImportServers() function from server-controller.go file")
}

func (controller *serverController) ExportServers() {
	fmt.Println("Execute ExportServers() function from server-controller.go file")
}

func (controller *serverController) CheckServerExistence(ip string) bool {
	return controller.serverService.CheckServerExistence(ip)
}

func (controller *serverController) CheckServerName(name string) bool {
	return controller.serverService.CheckServerName(name)
}

func (controller *serverController) ReportServerStatus() {
	// Starting with elastic.v5, you must pass a context to execute each service
	ctx := context.Background()

	// Take all servers from database
	servers_for_report, err := serverDatabase.ViewServers(0, 0, 0, "", "", "")
	if err != nil {
		panic(err)
	}

	// log.Println(servers_for_report)
	// Create a xlsx file contain information about servers
	new_file := excelize.NewFile()
	defer func() {
		if err := new_file.Close(); err != nil {
			panic(err)
		}
	}()

	new_file.SetCellValue("Sheet1", "A1", "ID")
	new_file.SetCellValue("Sheet1", "B1", "Server Name")
	new_file.SetCellValue("Sheet1", "C1", "IPv4")
	new_file.SetCellValue("Sheet1", "D1", "User")
	new_file.SetCellValue("Sheet1", "E1", "Password")
	new_file.SetCellValue("Sheet1", "F1", "Status")
	new_file.SetCellValue("Sheet1", "G1", "Created At")
	new_file.SetCellValue("Sheet1", "H1", "Updated At")
	new_file.SetCellValue("Sheet1", "I1", "Percentage Of Uptime")
	new_file.SetCellValue("Sheet1", "J1", "Downtime Intervals")
	for idx, server := range servers_for_report {
		_ = server
		row_idx := strconv.FormatInt(int64(idx+2), 10)
		cell := "A" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.Id)
		cell = "B" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.Name)
		cell = "C" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.Ipv4)
		cell = "D" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.User)
		cell = "E" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.Password)
		cell = "F" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.Status)
		cell = "G" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.CreatedAt)
		cell = "H" + row_idx
		new_file.SetCellValue("Sheet1", cell, server.UpdatedAt)
	}

	// Calculate the uptime status for each
	// - Establish connection to elasticsearch
	var client *elastic.Client
	err = try.Do(func(attempt int) (bool, error) {
		var err error
		client, err = elastic.NewClient(elastic.SetURL("http://elasticsearch:9200"))
		if err != nil {
			time.Sleep(10 * time.Second) // wait a minute
		}
		return attempt < 5, err
	})
	if err != nil {
		panic(err)
	}

	// Use the IndexExists service to check if a specified index exists.
	exists, err := client.IndexExists("myindex").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
	if !exists {
		// Create a new index.
		createIndex, err := client.CreateIndex("myindex").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
		if !createIndex.Acknowledged {
			// Not acknowledged
		}
	}

	// - Loop through servers, calculate the uptime status and save to xlsx file
	for idx, server := range servers_for_report {
		// Take ip of each server
		ip := server.Ipv4
		// Search records with a term query - ip
		termQuery := elastic.NewTermQuery("ipv4", ip)
		searchResult, err := client.Search().
			Index("myindex").         // search in index "myindex"
			Type("mytype").           // type "mytype"
			Query(termQuery).         // specify the query
			Sort("timestamp", false). // sort by "timestamp" field, ascending
			From(0).Size(100).        // take documents 0-9
			Pretty(true).             // pretty print request and response JSON
			Do(ctx)                   // execute

		if err != nil {
			// Handle error
			panic(err)
		}
		var down_time_records []time_event_element
		var status_stack []string

		var lasttime_to_subtract time.Time
		current_time := time.Now()
		// log.Println("Current time:", current_time)
		threehoursago_time := current_time.Add(-3 * time.Hour)
		// log.Println("3 Hours Ago time:", threehoursago_time)
		server_created_time := server.CreatedAt
		// log.Println("Server Created time:", server_created_time)
		// Check if three hours ago time point is newer or older than the time this server is added to sms
		if server_created_time.Before(threehoursago_time) {
			lasttime_to_subtract = threehoursago_time
		} else {
			lasttime_to_subtract = *server_created_time
		}

		// log.Println("Last time to subtract:", lasttime_to_subtract)

		minuend := current_time
		var serverHBR ServerHeartBeatResponse
		for _, item := range searchResult.Each(reflect.TypeOf(serverHBR)) {
			if s, ok := item.(ServerHeartBeatResponse); ok {
				if s.Time.Before(lasttime_to_subtract) {
					break
				}
				// fmt.Printf("Server %s: %s\n", s.Ipv4, s.Time)
				subtrahend := s.Time
				if minuend.Sub(subtrahend) > 2*time.Minute {
					down_minutes := int64(minuend.Sub(subtrahend) / time.Minute)
					for i := 1; i < int(down_minutes); i++ {
						status_stack = append(status_stack, "Down")
					}
					down_time_records = append(down_time_records, time_event_element{
						from:     subtrahend,
						to:       minuend,
						duration: minuend.Sub(subtrahend),
					})
				} else {
					status_stack = append(status_stack, "Up")
				}
				minuend = subtrahend
			}
		}
		if minuend.Sub(lasttime_to_subtract) > 2*time.Minute && lasttime_to_subtract.Before(minuend) {
			down_minutes := int64(minuend.Sub(lasttime_to_subtract) / time.Minute)
			for i := 1; i < int(down_minutes); i++ {
				status_stack = append(status_stack, "Down")
			}
			down_time_records = append(down_time_records, time_event_element{
				from:     lasttime_to_subtract,
				to:       minuend,
				duration: minuend.Sub(lasttime_to_subtract),
			})
		} else {
			status_stack = append(status_stack, "Up")
		}

		// fmt.Println(status_stack)
		//Create a dictionary of values for each element
		dict := make(map[string]int)
		for _, status := range status_stack {
			dict[status] = dict[status] + 1
		}
		// fmt.Println(dict)
		uptime_percentage := float64(dict["Up"]) / float64(dict["Up"]+dict["Down"])
		uptime_percentage = uptime_percentage * 100
		uptime_percentage_str := fmt.Sprintf("%0.2f%%", uptime_percentage)
		// fmt.Println(uptime_percentage_str)
		// for _, e := range down_time_records {
		// 	fmt.Println(e.from)
		// 	fmt.Println(e.to)
		// 	fmt.Println(e.duration)
		// }

		down_time_intervals_value := "["
		for _, e := range down_time_records {
			down_time_intervals_value += "{"
			down_time_intervals_value += fmt.Sprint("from:", e.from)
			down_time_intervals_value += fmt.Sprint("from:", e.to)
			down_time_intervals_value += fmt.Sprint("from:", e.duration)
			down_time_intervals_value += "},"
		}
		down_time_intervals_value += "]"
		// Add uptime status to each record in xlsx
		row_idx := strconv.FormatInt(int64(idx+2), 10)
		cell := "I" + row_idx
		new_file.SetCellValue("Sheet1", cell, uptime_percentage_str)
		cell = "J" + row_idx
		new_file.SetCellValue("Sheet1", cell, down_time_intervals_value)

		if err := new_file.SaveAs("./data-files/reported-servers.xlsx"); err != nil {
			panic(err)
		}
	}
	// Send file to admin each 30 minutes

	m := mail.NewMessage()

	m.SetHeader("From", "nguyenminhmanh060301@gmail.com")

	m.SetHeader("To", "nguyenminhmannh2001@gmail.com")

	m.SetHeader("Subject", "Report Server Status")

	m.SetBody("text/html", "Hello <b>Admin</b>! <br> This is the report email.")

	m.Attach("./data-files/reported-servers.xlsx")

	d := mail.NewDialer("smtp.gmail.com", 587, "nguyenminhmanh060301@gmail.com", "qxtmjeozmltamikz")

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

}

func (controller *serverController) AutoReportServerStatus() {
	for {
		log.Println("Send email to admin every 30 minutes")
		controller.ReportServerStatus()
		time.Sleep(30 * time.Minute)
	}
}
