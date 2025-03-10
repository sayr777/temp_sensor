package main

import (
    "encoding/binary"
    "machine"
    "math"
    "time"

    "tinygo.org/x/bluetooth"
    "tinygo.org/x/drivers/ds18b20"
    "tinygo.org/x/drivers/onewire"
)

var (
    adapter  = bluetooth.DefaultAdapter
    sensor   ds18b20.Device
    tempChar *bluetooth.Characteristic
)

// UUID сервиса и характеристики
var (
    serviceUUID = bluetooth.NewUUID([16]byte{0x21, 0x9C, 0x45, 0xCE, 0x11, 0xE9, 0x89, 0x5B, 0x78, 0x4A, 0x92, 0xDF, 0x6B, 0xBC, 0xD7, 0x2A})
    charUUID    = bluetooth.NewUUID([16]byte{0x21, 0x9C, 0x45, 0xCE, 0x11, 0xE9, 0x89, 0x5B, 0x78, 0x4A, 0x92, 0xDF, 0x6B, 0xBC, 0xD7, 0x2B})
)

func main() {
    // Инициализация датчика DS18B20
    pin := machine.P0_17
    ow := onewire.New(pin)
    sensor = ds18b20.New(ow)

    // Инициализация BLE
    must("enable BLE", adapter.Enable())

    // Настройка рекламы
    adv := adapter.DefaultAdvertisement()
    must("config adv", adv.Configure(bluetooth.AdvertisementOptions{
        LocalName:    "BLE Temp Sensor",
        ServiceUUIDs: []bluetooth.UUID{serviceUUID},
    }))
    must("start adv", adv.Start())

    // Создание BLE сервиса с характеристикой
    must("add service", adapter.AddService(&bluetooth.Service{
        UUID: serviceUUID,
        Characteristics: []bluetooth.CharacteristicConfig{
            {
                UUID:  charUUID,
                Flags: bluetooth.CharacteristicReadPermission | bluetooth.CharacteristicNotifyPermission,
                WriteEvent: func(client bluetooth.Connection, offset int, value []byte) {
                    // Для примера: обрабатываем запись в характеристику (можно добавить обработку)
                    println("Received data:", value)
                },
            },
        },
    }))

    // Основной цикл
    for {
        tempBytes, err := getTemperatureBytes()
        if err != nil {
            println("Error reading temp:", err.Error())
            time.Sleep(3 * time.Second)
            continue
        }

        sendTemperature(tempBytes)
        time.Sleep(3 * time.Second)
    }
}

func getTemperatureBytes() ([]byte, error) {
    sensor.RequestTemperature(nil) // Запрос температуры
    time.Sleep(750 * time.Millisecond) // Время преобразования температуры

    temp, err := sensor.ReadTemperature(nil)
    if err != nil {
        return nil, err
    }

    buf := make([]byte, 4)
    binary.LittleEndian.PutUint32(buf, math.Float32bits(float32(temp)))
    return buf, nil
}

func sendTemperature(tempBytes []byte) {
    if tempChar != nil {
        tempChar.Write(tempBytes) // Записываем данные в характеристику
    }
}

func must(action string, err error) {
    if err != nil {
        panic("failed to " + action + ": " + err.Error())
    }
}
