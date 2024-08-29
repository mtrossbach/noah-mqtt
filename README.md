# noah-mqtt
![License](https://img.shields.io/github/license/mtrossbach/noah-mqtt) ![GitHub last commit](https://img.shields.io/github/last-commit/mtrossbach/noah-mqtt) ![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/mtrossbach/noah-mqtt)

UNDER CONSTRUCTION!

`noah-mqtt` is a standalone application designed to retrieve data and metrics from your Growatt NOAH 2000 home battery used in balcony power plants. It publishes this information to an MQTT broker, making it easily accessible for Home Assistant or other applications.

The application features Home Assistant auto-discovery, allowing your NOAH devices to be automatically recognized and integrated with Home Assistant via the MQTT integration.


# ![HomeAssistant screenshot](/assets/ha-screenshot.png)

# Run the application standalone

## Option 1: Running `noah-mqtt` with Docker

To run the latest version of `noah-mqtt` using Docker, follow these steps:

1. **Install Docker**: Ensure Docker is installed on your system. You can download Docker Desktop from [Docker’s official website](https://www.docker.com/products/docker-desktop).

2. **Open a Terminal**:
   - **Windows**: Use Command Prompt or PowerShell.
   - **Linux/macOS**: Use the Terminal.

3. **Execute the Docker Command**: Run the following command, replacing the placeholders with your actual values:

   ```
   docker run -name noah-mqtt -e GROWATT_USERNAME=myusername -e GROWATT_PASSWORD=mypassword -e MQTT_HOST=localhost -e MQTT_PORT=1883 ghcr.io/mtrossbach/noah-mqtt:latest
   ```
   
- Replace myusername with your Growatt username.
- Replace mypassword with your Growatt password.
- Replace localhost with the hostname or IP address of your MQTT broker.
- Replace 1883 with the port number your MQTT broker uses (default is 1883).

The application will connect to your MQTT broker and retrieve all metrics and data for your NOAH devices.

## Option 2: Downloading and running a prebuilt binary

If you prefer not to compile the binary yourself, you can download a prebuilt version:

1. **Download the Binary**: Go to the [Releases](https://github.com/mtrossbach/noah-mqtt/releases) page of the repository and download the prebuilt binary for your operating system and system architecture.

2. **Extract the Binary**: If the binary is compressed (e.g., in a zip or tar file), extract it to a directory of your choice.

3. **Run the Application**: Open a terminal in the directory containing the binary and run it using the appropriate command for your OS, setting the necessary environment variables:

   - **Windows** (Command Prompt):

     ```sh
     set GROWATT_USERNAME=myusername
     set GROWATT_PASSWORD=mypassword
     set MQTT_HOST=localhost
     set MQTT_PORT=1883
     noah-mqtt.exe
     ```

   - **Windows** (PowerShell):

     ```sh
     $env:GROWATT_USERNAME=„myusername“
     $env:GROWATT_PASSWORD=„mypassword“
     $env:MQTT_HOST=„localhost“
     $env:MQTT_PORT=„1883“
     .\noah-mqtt.exe
     ```

   - **Linux/macOS**:

     ```sh
     GROWATT_USERNAME=myusername GROWATT_PASSWORD=mypassword MQTT_HOST=localhost MQTT_PORT=1883 ./noah-mqtt
     ```

Again, replace `myusername`, `mypassword`, `localhost`, and `1883` with your actual Growatt account details and MQTT broker information.

## Option 3: Compiling the binary yourself

To compile the binary yourself, ensure you have Go installed on your machine:

1. **Install Go**: Download and install the latest version of Go from [the official Go website](https://golang.org/dl/).

2. **Clone the Repository**: Open a terminal and run the following command to clone the repository:
        
        git clone https://github.com/mtrossbach/noah-mqtt.git
        cd noah-mqtt

3. **Build the application**:

        go build -o noah-mqtt cmd/noah-mqtt/main.go

Afterwards follow the instructions for running the application from option 2.


# Integration into HomeAssistant

`noah-mqtt` interacts with Home Assistant by publishing data from your Growatt NOAH 2000 home battery to an MQTT broker. This setup allows Home Assistant to subscribe to and integrate this data seamlessly into its ecosystem.

![Home Assistant Integration](./assets/noah-mqtt-ha-dark.drawio.png#gh-dark-mode-only)
![Home Assistant Integration](./assets/noah-mqtt-ha.drawio.png#gh-light-mode-only)

If you’re already using MQTT with other integrations like zigbee2mqtt or AhoyDTU, you already have the MQTT integration configured and active. In this case, you can skip step 1 and 2 as your existing setup should work with `noah-mqtt`.

The following integration process for `noah-mqtt` with Home Assistant works for all installation methods, regardless of how Home Assistant is installed—whether it’s through Home Assistant OS, Home Assistant Supervised, or Home Assistant Container. 

1. **Set Up an MQTT Broker**:  
   Ensure you have an MQTT broker running, such as [Mosquitto](https://mosquitto.org/), and that it’s accessible from both Noah-mqtt and Home Assistant.

2. **Check MQTT Integration in Home Assistant**:  
   - Navigate to **Settings** > **Devices & Services** in Home Assistant.
   - Click **Add Integration** and select „MQTT“.
   - Enter your MQTT broker details (hostname, port, username, password).
   - Test the connection to ensure it’s working correctly.

3. **Run noah-mqtt**:  
   Start `noah-mqtt` using the appropriate configuration for your MQTT broker.

4. **Verify Device Discovery**:  
   Check **Devices** and **Entities** under **Settings** > **Devices & Services** in Home Assistant to confirm that your Noah devices are automatically discovered.

By following these steps, `noah-mqtt` will communicate with Home Assistant via your MQTT broker, also supporting automatic device discovery. If you already have MQTT set up, it should integrate seamlessly with your existing configuration.

# Configuration

You can configure `noah-mqtt` using the following environment variables:

| Environment Variable                              | Description                                                                             | Default       |
|:--------------------------------------------------|:----------------------------------------------------------------------------------------|:--------------| 
| `LOG_LEVEL`                                       | Sets the logging level of the application                                               | INFO          |
| `POLLING_INTERVAL`                                | Time in seconds between fetching new status data                                        | 10            |
| `BATTERY_DETAILS_POLLING_INTERVAL`                | Time in seconds between fetching battery details (per battery SoC & temperature).       | 60            |
| `PARAMETER_POLLING_INTERVAL`                      | Time in seconds between fetching parameter data (system-output-power, charging limits). | 60            |
| `GROWATT_USERNAME`                                | Your Growatt account username (required)                                                | -             |
| `GROWATT_PASSWORD`                                | Your Growatt account password (required)                                                | -             |
| `MQTT_HOST`                                       | Address of your MQTT broker (required)                                                  | -             |
| `MQTT_PORT`                                       | Port number of your MQTT broker                                                         | 1883          |
| `MQTT_CLIENT_ID`                                  | Identifier for the MQTT client                                                          | noah-mqtt     |
| `MQTT_USERNAME`                                   | Username for connecting to your MQTT broker                                             | -             |
| `MQTT_PASSWORD`                                   | Password for connecting to your MQTT broker                                             | -             |
| `MQTT_TOPIC_PREFIX`                               | Prefix for MQTT topics used by Noah-mqtt                                                | noah2mqtt     |
| `HOMEASSISTANT_TOPIC_PREFIX`                      | Prefix for topics used by Home Assistant                                                | homeassistant |

Adjust these settings to fit your environment and requirements.


# Run as Home Assistant add-on (Home Assistant OS, Home Assistant Supervised)

If you are using Home Assistant OS or Home Assistant Supervised you can run `noah-mqtt` as a Home Assistant add-on, which provides seamless integration with your Home Assistant setup.
This option leverages the add-on system to manage and run `noah-mqtt` directly on your Home Assistant instance.

#### Steps to Use the Home Assistant Add-on

1. **Add the Repository:**
   - Open your Home Assistant web interface.
   - Navigate to **Settings** > **Add-ons** > **Add-on Store**.
   - Click on the three-dot menu in the top right corner and select **Repositories**.
   - Add the following URL: `https://github.com/mtrossbach/hassio-addons`.

[![Open your Home Assistant instance and show the add add-on repository dialog with a specific repository URL pre-filled.](https://my.home-assistant.io/badges/supervisor_add_addon_repository.svg)](https://my.home-assistant.io/redirect/supervisor_add_addon_repository/?repository_url=https%3A%2F%2Fgithub.com%2Fmtrossbach%2Fhassio-addons)

2. **Install the Add-on:**
   - Search for the `noah-mqtt` add-on within the Add-on Store.
   - Click on the add-on and select **Install**.

3. **Configure the Add-on:**
   - After installation, configure the add-on settings by providing your **Growatt username** and **Growatt password** and setup the other options as needed.

4. **Start the Add-on:**
   - Click **Start** to launch the `noah-mqtt` add-on.

The Home Assistant add-on provides an easy and integrated way to run `noah-mqtt`, allowing you to manage it directly from the Home Assistant interface.

For more detailed information and updates, visit the [repository](https://github.com/mtrossbach/hassio-addons).
