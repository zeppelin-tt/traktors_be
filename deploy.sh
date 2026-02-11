#!/bin/bash
set -e

# =============================================================================
# Деплой traktors_be на Ubuntu VM
#
# Использование:
#   ./deploy.sh <host> [user] [путь_к_ssh_ключу]
#
# Примеры:
#   ./deploy.sh 185.10.20.30 ubuntu ~/.ssh/id_rsa
#   ./deploy.sh 185.10.20.30                        # user=ubuntu, ключ по умолчанию
# =============================================================================

HOST="${1:?Укажите host: ./deploy.sh <host> [user] [ssh_key]}"
USER="${2:-ubuntu}"
KEY="${3:-}"

APP_DIR="/opt/traktors_be"
GO_VERSION="1.22.6"

# --- helpers -----------------------------------------------------------------
SSH_ARGS=(-o StrictHostKeyChecking=no -o BatchMode=yes)
[[ -n "$KEY" ]] && SSH_ARGS+=(-i "$KEY")

run()    { ssh "${SSH_ARGS[@]}" "${USER}@${HOST}" "$@"; }
upload() { scp "${SSH_ARGS[@]}" "$@" "${USER}@${HOST}:${APP_DIR}/"; }

# --- connectivity ------------------------------------------------------------
echo "→ Подключение к ${USER}@${HOST}…"
run echo "  OK"

# --- Go ----------------------------------------------------------------------
echo "→ Go…"
if run "which go" >/dev/null 2>&1; then
    echo "  уже установлен"
else
    echo "  устанавливаю ${GO_VERSION}…"
    run "curl -fsSL https://go.dev/dl/go${GO_VERSION}.linux-amd64.tar.gz | sudo tar -C /usr/local -xz"
    run "echo 'export PATH=\$PATH:/usr/local/go/bin' | sudo tee /etc/profile.d/go.sh >/dev/null"
fi

# --- MongoDB -----------------------------------------------------------------
echo "→ MongoDB…"
if run "which mongod" >/dev/null 2>&1; then
    echo "  уже установлен"
else
    echo "  устанавливаю…"
    run 'curl -fsSL https://www.mongodb.org/static/pgp/server-7.0.asc | sudo gpg --yes --dearmor -o /usr/share/keyrings/mongodb-server-7.0.gpg'
    run 'echo "deb [ signed-by=/usr/share/keyrings/mongodb-server-7.0.gpg] https://repo.mongodb.org/apt/ubuntu jammy/mongodb-org/7.0 multiverse" | sudo tee /etc/apt/sources.list.d/mongodb-org-7.0.list >/dev/null'
    run 'sudo apt-get update -qq'
    run 'sudo apt-get install -y -qq mongodb-org'
fi

# --- upload & build ----------------------------------------------------------
echo "→ Копирование файлов…"
run "sudo mkdir -p ${APP_DIR} && sudo chown ${USER}:${USER} ${APP_DIR}"
run "sudo mkdir -p ${APP_DIR}/uploads && sudo chown ${USER}:${USER} ${APP_DIR}/uploads"
upload main.go model.go handlers.go media_handlers.go go.mod go.sum

echo "→ Сборка…"
run "export PATH=\$PATH:/usr/local/go/bin && cd ${APP_DIR} && go mod tidy && go build -o traktors_be ."

# --- systemd -----------------------------------------------------------------
echo "→ Настройка systemd…"
cat << EOF | run "sudo tee /etc/systemd/system/traktors_be.service >/dev/null"
[Unit]
Description=Traktors Backend API
After=network.target mongod.service

[Service]
Type=simple
User=${USER}
Group=${USER}
WorkingDirectory=${APP_DIR}
ExecStart=${APP_DIR}/traktors_be
Restart=always
RestartSec=5
Environment=PORT=8080
Environment=MONGO_URI=mongodb://localhost:27017
Environment=DB_NAME=traktors
Environment=UPLOAD_DIR=${APP_DIR}/uploads
Environment=BASE_URL=http://${HOST}:8080

[Install]
WantedBy=multi-user.target
EOF

echo "→ Старт сервисов…"
run "sudo systemctl daemon-reload"
run "sudo systemctl enable --now mongod"
run "sudo systemctl enable --now traktors_be"

echo ""
run "sudo systemctl status traktors_be --no-pager" || true
echo ""
echo "✓ Готово.  API:  http://${HOST}:8080/tractors"
