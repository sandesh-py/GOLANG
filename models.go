package main

type Car struct {
    ID        string `json:"id"`
    Brand     string `json:"brand"`
    Number    string `json:"number"`
    Type      string `json:"type"`
    Incoming  string `json:"incoming_time"`
    Outgoing  string `json:"outgoing_time"`
    Slot      string `json:"parking_slot"`
}