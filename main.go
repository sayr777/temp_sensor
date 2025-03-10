package main

import (
    "machine"
    "time"

    "tinygo.org/x/drivers/onewire"
    "tinygo.org/x/drivers/ds18b20"
)

func main() {
    // Определяем пин для DS18B20
    pin := machine.P0_17

    // Инициализируем OneWire шину
    ow := onewire.New(pin)
    sensor := ds18b20.New(ow)

    for {
        time.Sleep(3 * time.Second)

        println()
        println("Device:", machine.Device)

        println()
        println("Request Temperature.")

        // Запрашиваем температуру, используя SKIP_ROM (так как датчик один)
        sensor.RequestTemperature(nil) // nil означает SKIP_ROM

        // Ждём время преобразования температуры (увеличили до 2 секунд)
        time.Sleep(2 * time.Second)

        println()
        println("Read Temperature")

        // Читаем RAW данные
        raw, err := sensor.ReadTemperatureRaw(nil)
        if err != nil {
            println("Failed to read raw temperature:", err)
            continue
        }
        println("Raw temperature data:", raw)

        // Читаем температуру
        t, err := sensor.ReadTemperature(nil)
        if err != nil {
            println("Failed to read temperature:", err)
            continue
        }

        // Выводим температуру в миллиградусах Цельсия
        println("Temperature in celsius milli degrees (°C/1000):", t)
    }
}