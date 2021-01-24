module example.com/main

go 1.15

replace example.com/kconnect => ./kconnect

replace example.com/datasorter => ./datasorter

replace example.com/datagen => ./datagen

replace example.com/logging => ./logging

require (
	example.com/datagen v0.0.0-00010101000000-000000000000
	example.com/datasorter v0.0.0-00010101000000-000000000000
	example.com/kconnect v0.0.0-00010101000000-000000000000
	example.com/logging v0.0.0-00010101000000-000000000000
)
