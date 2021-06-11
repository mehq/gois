package main

func main() {
	options := &Options{
		query:    "cats",
		explicit: true,
		gif:      false,
		gray:     false,
		height:   0,
		width:    0,
	}
	ScrapeBing(options)
}
