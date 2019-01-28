package main

type exporter interface {
	setup()
	export(m []metrics) error
}
