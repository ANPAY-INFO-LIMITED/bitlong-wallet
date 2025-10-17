# Function to check if the previous command executed successfully
check_error() {
    if [ $? -ne 0 ]; then
        echo "Error: $1" | tee -a /var/log/icn_script_errors.log
        echo "Press Enter to return to the menu..."
        read -r < /dev/tty
        return 1
    fi
}

# Check if the script is running as root
if [ "$EUID" -ne 0 ]; then
    echo "Error: This script must be run as root"
    exit 1
fi

# litd update function
update_litd() {
    if [ ! -f /usr/local/bin/litd ]; then
        echo "Error: litd binary not found in /usr/local/bin"
        echo "Please run option 10 to install litd first."
        echo "Press Enter to return to the menu..."
        read -r < /dev/tty
        return 1
    fi

    echo "Downloading SHA1 hash file..."
    curl -o /tmp/litd.sha1 --user 1IzF:G9Me https://api.btc.microlinktoken.com:28173/f/litd.sha1
    check_error "Failed to download SHA1 hash file" || return 1

    current_sha1=$(sha1sum /usr/local/bin/litd | awk '{print $1}')
    check_error "Failed to calculate SHA1 for current litd binary" || return 1

    downloaded_sha1=$(cat /tmp/litd.sha1)
    check_error "Failed to read downloaded SHA1 file" || return 1

    if [ "$current_sha1" = "$downloaded_sha1" ]; then
        echo "litd is already up to date"
        rm /tmp/litd.sha1
        check_error "Failed to delete litd.sha1 file" || return 1
        echo "Press Enter to return to the menu..."
        read -r < /dev/tty
        return 0
    fi

    echo "New version detected, updating litd..."

    systemctl stop litd.service
    check_error "Failed to stop litd service" || return 1

    echo "Downloading new bin.tar.gz..."
    max_retries=3
    retry_count=0
    until [ $retry_count -ge $max_retries ]; do
        curl -o /tmp/bin.tar.gz --user 1IzF:G9Me https://api.btc.microlinktoken.com:28173/bin.tar.gz && break
        retry_count=$((retry_count + 1))
        echo "Retry $retry_count/$max_retries ..."
        sleep 1
    done
    if [ $retry_count -ge $max_retries ]; then
        echo "Error: Failed to download bin.tar.gz after $max_retries attempts"
        echo "Press Enter to return to the menu..."
        read -r < /dev/tty
        return 1
    fi

    if [ ! -f "/tmp/bin.tar.gz" ] || [ ! -s "/tmp/bin.tar.gz" ]; then
        echo "Error: bin.tar.gz file is missing or empty"
        echo "Press Enter to return to the menu..."
        read -r < /dev/tty
        return 1
    fi

    echo "Extracting bin.tar.gz..."
    tar -zxvf /tmp/bin.tar.gz -C /tmp/
    check_error "Failed to extract bin.tar.gz" || return 1

    BIN_DIR="/tmp/bin"
    if [ ! -d "$BIN_DIR" ]; then
        echo "Error: bin directory not found in /tmp"
        echo "Press Enter to return to the menu..."
        read -r < /dev/tty
        return 1
    fi

    echo "Copying all files from $BIN_DIR to /usr/local/bin..."
    cp -r "$BIN_DIR"/* /usr/local/bin/
    check_error "Failed to copy files to /usr/local/bin" || return 1

    echo "Cleaning up $BIN_DIR and bin.tar.gz..."
    rm -rf "$BIN_DIR"
    check_error "Failed to delete $BIN_DIR" || return 1
    rm -f /tmp/bin.tar.gz
    check_error "Failed to delete bin.tar.gz" || return 1
    rm -f /tmp/litd.sha1
    check_error "Failed to delete litd.sha1 file" || return 1

    systemctl start litd.service
    check_error "Failed to start litd service" || return 1

    echo "Checking updated litd service status..."
    systemctl status litd.service --no-pager
    check_error "Failed to check updated litd service status" || return 1

    echo "litd successfully updated and restarted"
}

update_litd
