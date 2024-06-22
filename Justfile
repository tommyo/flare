default:
    @just --list

# generate protobuf
gen:
    @buf lint
    @buf generate