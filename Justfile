set dotenv-load

cwd := justfile_directory()

make-model:
    gojsonschema -p stms -o model.go run.schema.json

install-gojsonschema:
    git clone https://github.com/omissis/go-jsonschema
    cd go-jsonschema
    go install cmd/gojsonschema

[unix]
sqlc +CMD:
    docker run --rm -v "{{cwd}}:/src" -w /src kjconroy/sqlc {{CMD}}

[windows]
sqlc +CMD:
    #!powershell
    docker run --rm -v "{{cwd}}:/src" -w /src kjconroy/sqlc {{CMD}}

sqlc-gen: (sqlc "-f" "/src/sql/sqlc.yaml" "generate")
    

serve:
    go run cmd/stsms/main.go

migrate DEST:
    cd sql; tern migrate --destination {{DEST}}

upload-runs DIR URL:
    cd "{{DIR}}" && for fn in $(ls); do echo -n "$fn "; curl -X POST "{{URL}}" --data-binary "@$fn"; echo; done
