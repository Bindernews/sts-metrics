set dotenv-load
win_prefix := if os() == "windows" { 'MSYS2_ARG_CONV_EXCL="*"' } else { '' }
py_cmd := if os() == "windows" { 'python' } else { 'python3' }

cwd := justfile_directory()

# Build the JSON model using gojsonschema, and then inserts additional code because
# that only gets us part-way.
make-model:
    gojsonschema -p web -o web/json_model.go run.schema.json
    # Insert imports, change PlayId to be uuid.UUID, add Extra field and code
    # to process said extra field
    sed -ri \
        -e 's/PlayId string/PlayId uuid.UUID/g' \
        -e '/^import "fmt"/a import "github.com/google/uuid"\nimport "github.com/samber/lo"' \
        -e '/WinRate float64/a\	// Additional fields\n	Extra map[string]any `json:"-"`' \
        -e '/\*j = RunSchema/i\	plain.Extra = lo.OmitByKeys(raw, runSchemaJsonKeys)' \
        web/json_model.go
    jq '.properties|keys[]' run.schema.json |\
        awk 'BEGIN{ print "\nvar runSchemaJsonKeys = []string{" } { print "\t"$0"," } END { print "}" }' \
        >>web/json_model.go

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

# Generate a secret suitable for securing cookies
make-secret:
    {{py_cmd}} -c 'import os,base64; print(base64.urlsafe_b64encode(os.urandom(32)).decode())'
    
serve:
    go run cmd/stsms/main.go

# Use tern to perform migrations
migrate DEST:
    cd sql && tern migrate --destination {{DEST}}

# Migrate to 0, then to last/latest
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
