set dotenv-load
win_prefix := if os() == "windows" { 'MSYS2_ARG_CONV_EXCL="*"' } else { '' }

cwd := justfile_directory()

make-model:
    gojsonschema -p stms -o json_model.go run.schema.json
    sed -ri \
        -e 's/PlayId string/PlayId uuid.UUID/g;/^import "fmt"/a import "github.com/google/uuid"' \
        -e '/WinRate float64/a\	// Additional fields\n	Extra map[string]any' \
        -e '/\*j = RunSchema/i\	plain.Extra = extProcessRaw(raw)' \
        json_model.go

# Outputs a list of known fields in RunSchemaJson.
# This can be used to update 'extProcessRaw' in json_to_orm.go
list-schema-fields:
    jq '.properties|keys[]' run.schema.json

install-gojsonschema DEST:
    git clone https://github.com/omissis/go-jsonschema {{DEST}}/go-jsonschema
    cd {{DEST}}/go-jsonschema && go install cmd/gojsonschema

# Installs the "tern" command via "go install"
install-tern:
    go install github.com/jackc/tern@latest

# Run a command in the sqlc container
sqlc +CMD:
    {{win_prefix}} docker run --rm -v "{{cwd}}:/src" -w /src kjconroy/sqlc:1.17.2 {{CMD}}

# Generate the ORM package
sqlc-gen: (sqlc "-f" "/src/sql/sqlc.yaml" "generate")
    
# Generate a new self-signed certificate in ./data
# On Windows it's best to run this within WSL
make-cert CERT_NAME="data/ssl":
    openssl req -new -subj "/C=US/ST=Ohio/CN=localhost" \
      -newkey rsa:2048 -nodes -keyout "{{CERT_NAME}}.key" -out "{{CERT_NAME}}.pem"
    openssl x509 -req -days 365 -in "{{CERT_NAME}}.pem" -signkey "{{CERT_NAME}}.key" -out "{{CERT_NAME}}.crt"

serve:
    go run cmd/stsms/main.go

migrate DEST:
    cd sql && tern migrate --destination {{DEST}}

reset-db:
    just migrate 0
    just migrate last

upload-runs DIR URL:
    cd "{{DIR}}" && \
    for fn in $(find . -iname "*.run" -o -iname "*.run.gz"); do \
      echo -n "$fn "; \
      if [ "${fn##*.}" == "gz" ]; then gzip -dc "$fn"; else cat "$fn"; fi | curl -X POST "{{URL}}" --data-binary "@-"; \
      echo; \
    done
