MODULE_NAME="nagare_computer_use"
MODULE_DIR="$HOME/.nagare/modules"
cargo build
install -Dm755 "../../target/debug/$MODULE_NAME" "$MODULE_DIR/$MODULE_NAME"
