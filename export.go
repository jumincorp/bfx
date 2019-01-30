package main

type exporter interface {
	setup()
	export([]metric) error
}
