
update-assets:
	(cd web; gulp)
	mkdir -p pkg/web/static
	esc -o pkg/web/static/static.go -pkg static -prefix web/dist web/dist
