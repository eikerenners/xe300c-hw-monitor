package main

import (
	"encoding/json"
	"fmt"
)

type Network struct {
	Online    bool   `json:"online"`
	Up        bool   `json:"up"`
	Interface string `json:"interface"`
}

type Wifi struct {
	Guest   bool   `json:"guest"`
	SSID    string `json:"ssid"`
	Up      bool   `json:"up"`
	Channel int    `json:"channel"`
	Band    string `json:"band"`
	Name    string `json:"name"`
	Passwd  string `json:"passwd"`
}

type Service struct {
	Status  int    `json:"status"`
	PeerID  int    `json:"peer_id,omitempty"`
	Name    string `json:"name"`
	GroupID int    `json:"group_id,omitempty"`
}

type Client struct {
	CableTotal    int `json:"cable_total"`
	WirelessTotal int `json:"wireless_total"`
}

type MCU struct {
	ChargeCnt      int     `json:"charge_cnt"`
	Temperature    float64 `json:"temperature"`
	ChargePercent  int     `json:"charge_percent"`
	ChargingStatus int     `json:"charging_status"`
}

type System struct {
	NetnatEnabled   bool      `json:"netnat_enabled"`
	GuestIP         string    `json:"guest_ip"`
	FlashApp        int       `json:"flash_app"`
	IPv6Enabled     bool      `json:"ipv6_enabled"`
	GuestNetmask    string    `json:"guest_netmask"`
	FlashFree       int       `json:"flash_free"`
	LoadAverage     []float64 `json:"load_average"`
	Mode            int       `json:"mode"`
	TZOffset        string    `json:"tzoffset"`
	LanNetmask      string    `json:"lan_netmask"`
	FlashTotal      int       `json:"flash_total"`
	MemoryTotal     int       `json:"memory_total"`
	MemoryFree      int       `json:"memory_free"`
	DDNSEnabled     bool      `json:"ddns_enabled"`
	Uptime          float64   `json:"uptime"`
	LanIP           string    `json:"lan_ip"`
	Timestamp       int       `json:"timestamp"`
	MCU             MCU       `json:"mcu"`
	MemoryBuffCache int       `json:"memory_buff_cache"`
}

type Result struct {
	Network []Network `json:"network"`
	Wifi    []Wifi    `json:"wifi"`
	Service []Service `json:"service"`
	Client  []Client  `json:"client"`
	System  System    `json:"system"`
}

type glStatusResponse struct {
	ID      int    `json:"id"`
	Jsonrpc string `json:"jsonrpc"`
	Result  Result `json:"result"`
}

func ParseGetStatusMsg(input string) (*glStatusResponse, error) {
	var response glStatusResponse
	err := json.Unmarshal([]byte(input), &response)
	if err != nil {
		return &glStatusResponse{}, err
	}
	return &response, nil
}

func TestGetStatusResponse() {
	jsonStr := `{"id":1,"jsonrpc":"2.0","result":{"network":[{"online":false,"up":false,"interface":"wan"},{"online":false,"up":false,"interface":"wwan"},{"online":false,"up":false,"interface":"tethering"},{"online":false,"up":false,"interface":"wan6"},{"online":false,"up":false,"interface":"wwan6"},{"online":false,"up":false,"interface":"tethering6"},{"online":true,"up":true,"interface":"modem_1_1_2"},{"online":false,"up":false,"interface":"modem_1_1_2_6"}],"wifi":[{"guest":false,"ssid":"REMOVED","up":true,"channel":0,"band":"2G","name":"default_radio0","passwd":"REMOVED"},{"guest":true,"ssid":"GL-XE300-efd-Guest","up":false,"channel":0,"band":"2G","name":"guest2g","passwd":"goodlife"}],"service":[{"status":1,"peer_id":4428,"name":"wgclient","group_id":7687},{"name":"wgserver","status":0},{"name":"ovpnclient","status":0},{"name":"ovpnserver","status":0}],"client":[{"cable_total":0,"wireless_total":2}],"system":{"netnat_enabled":false,"guest_ip":"192.168.9.1","flash_app":3375104,"ipv6_enabled":false,"guest_netmask":"255.255.255.0","flash_free":94781440,"load_average":[0.62,0.75,0.89],"mode":0,"tzoffset":"+0100","lan_netmask":"255.255.255.0","flash_total":134217728,"memory_total":124096512,"memory_free":48201728,"ddns_enabled":false,"uptime":20890.01,"lan_ip":"192.168.8.1","timestamp":1741510126,"mcu":{"charge_cnt":4,"temperature":38.1,"charge_percent":100,"charging_status":1},"memory_buff_cache":32677888}}}`

	response, err := ParseGetStatusMsg(jsonStr)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
	} else {
		fmt.Printf("Parsed Response: %+v\n", response)
	}
}
