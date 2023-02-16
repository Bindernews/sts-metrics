set windows-shell := ["powershell"]
set dotenv-load

cwd := justfile_directory()

make-model:
    gojsonschema -p stms -o model.go run.schema.json

install-gojsonschema:
    git clone https://github.com/omissis/go-jsonschema
    cd go-jsonschema
    go install cmd/gojsonschema

sqlc CMD:
    docker run --rm -v "{{cwd}}:/src" -w /src kjconroy/sqlc {{CMD}}
sqlc-gen:
    docker run --rm -v "{{cwd}}:/src" -w /src kjconroy/sqlc -f /src/sql/sqlc.yaml generate
    
serve:
    go run cmd/stsms/main.go

migrate DEST:
    cd sql; tern migrate -d {{DEST}}
