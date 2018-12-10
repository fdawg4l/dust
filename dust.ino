/*
 * Copyright 2018, fdawg4l@github.com
 * Permission is hereby granted, free of charge, to any person obtaining a copy of
 * this software and associated documentation files (the "Software"), to deal in
 * the Software without restriction, including without limitation the rights to
 * use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies
 * of the Software, and to permit persons to whom the Software is furnished to do
 * so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

#include <ESP8266WiFi.h>
#include <ESP8266WebServer.h>
#include <ESP8266mDNS.h>
#include <SoftwareSerial.h>
#include "PMS.h"
#include "DHT.h"

// Pin silkscreen vs reality
#define D2 4
#define D3 0
#define D4 2
#define D5 14

// Devices
#define LED D4
#define PMS_TX D2
#define PMS_RX D3
#define DHTPIN D5
#define DHTTYPE DHT11

// The fan takes 30s to stabilize and give useful readings.
#define PMS_FAN_DELAY 30000

#define SAMPLES 5

#define HOSTNAME "dust"

const char* ssid = "ShitakeMushrooms";
const char* password = "XXXXXXXXX";
char buffer[512] = "";

ESP8266WebServer server(80);
MDNSResponder mdns;

// RX - D2(GPIO4); Tx - D3(GPIO0)
// RX, TX
SoftwareSerial pmsSerial(PMS_TX,PMS_RX);
PMS pms(pmsSerial);
DHT dht(DHTPIN, DHTTYPE);

void setup() {
    // Poke the LED off
    pinMode(LED, OUTPUT);
    digitalWrite(LED, HIGH);

    Serial.begin(115200);

    // They say this thing supports i2c, but the pdf doesn't mention it.
    pmsSerial.begin(9600);

    WiFi.disconnect();
    WiFi.hostname(HOSTNAME);
    WiFi.begin(ssid, password);

    Serial.print("Attempting to connect to WPA SSID: ");
    Serial.println(ssid);

    // Wait for connection
    while (WiFi.status() != WL_CONNECTED) {
      delay(500);
      Serial.print(".");
    }

    // Dump the connection info
   Serial.println("");
   Serial.print("Connected to ");
   Serial.println(ssid);
   Serial.print("IP address: ");
   Serial.println(WiFi.localIP());

    // Start the MDNS
    if (mdns.begin(HOSTNAME)) {
       Serial.println("MDNS responder started");
    }

    // Add service to the MDNS
    mdns.addService("http", "tcp", 80);
    pms.passiveMode();
    pms.sleep();

    // DHT start
    dht.begin();

    server.on("/", []() {
      size_t len = SAMPLES;
      PMS::DATA data[len];
      char *buf = buffer;
      float hum = 0, temp = 0;
      uint idx = 0;

      if (!sample_temp(&temp, &hum)) {
         server.send(503, "text/plain", "oops\n");
         return;
      }

      // Gross hand rolling of json
      int n = sprintf(buf,
         "{\"t\": {"
                   "\"humidity_P\": %s,"
                    "\"temp_F\": %s},",
         String(hum).c_str(),
         String(temp).c_str());

      buf = buf + n;

      if (!sample_pms(data, len)) {
         server.send(503, "text/plain", "oops\n");
         return;
      }


      // start with an array
      n = sprintf(buf, "\"a\": [");
      buf += n;

      // add an object per item
      for (idx = 0; idx < len; idx++) {
        n = sprintf(buf,
          "{\"SP_1_0\": %d,"
           "\"SP_2_5\": %d,"
           "\"SP_10_0\": %d,"
           "\"AE_1_0\": %d,"
           "\"AE_2_5\": %d,"
           "\"AE_10_0\": %d}",
           data[idx].PM_SP_UG_1_0,
           data[idx].PM_SP_UG_2_5,
           data[idx].PM_SP_UG_10_0,
           data[idx].PM_AE_UG_1_0,
           data[idx].PM_AE_UG_2_5,
           data[idx].PM_AE_UG_10_0);

        buf = buf + n;

          // close out the array or append object
        if (idx + 1 == len) {
            sprintf(buf, "]}");
        } else {
            buf[0] = ',';
            buf++;
        }
      }

      server.send(200, "text/plain", buffer);
    });

    server.begin();
}

void loop() {
  server.handleClient();
}

// Sample will get a sample of data every 5s for the number of len.
// We turn on the light, wake up the sensor, wait 30s for the fan
// to do fan things, then start asking for data.  Then we turn the
// sensor off, and turn off the light.
bool sample_pms(PMS::DATA *out, size_t len) {
  uint8_t idx = 0;
  bool rv = true;
  digitalWrite(LED, LOW);
  pms.wakeUp();
  delay(PMS_FAN_DELAY);

  for(idx = 0; idx < len; idx++) {
    pms.requestRead();
    digitalWrite(LED, LOW);
    if (!pms.readUntil(out[idx])) {
      rv = false;
      goto out;
    }

    digitalWrite(LED, HIGH);
    delay(5000);
  }

out:

  pms.sleep();
  digitalWrite(LED, HIGH);

  return rv;
}

bool sample_temp(float *temp, float *hum) {
  uint8_t i = 0;
  bool rv = false;

  for (i = 0; i < 10; i ++) {
    *temp = dht.readTemperature(true);
    *hum = dht.readHumidity();

    if (!isnan(*hum) && !isnan(*temp)) {
      rv = true;
      break;
    }

    // Wait a few seconds between measurements
    delay(2000);
  }

  return rv;
}
