# SEMOVI


### Kafka example messages

Topics:
- hecho_transito
```bash
$ kafka-console-producer --broker-list PLAINTEXT://localhost:9092 --topic hecho_transito
> {"tiempo_hecho":1589279257759, "hecho_trans": 2, "placa": "ABC-123", "licencia": "C12345678", "distancia_viaje": "15-20", "tiempo_viaje": "100-200", "ubicación": "38.0088261,23.8042912"}
```
- stats_operador
```bash
$ kafka-console-producer --broker-list PLAINTEXT://localhost:9092 --topic stats_operador
> {"fecha_produccion":1589279257759,"id_operador":"2d2ec778-b89e-4db5-9628-123fd99f0b91","genero":1,"cant_viajes":29,"tiempo_registro":44,"edad":"28-32","horas_conectado":"9-17","horas_viaje":"0-24","ingreso_totales":"$0-999"}
```
- viajes_agregados
```bash
$ kafka-console-producer --broker-list PLAINTEXT://localhost:9092 --topic viajes_agregados
> {"fecha": 1589279257759,"id_proveedor": "test","tot_viajes": 12,"tot_veh_viaje": 11,"tot_veh_disp": 1,"dist_pasajero": 1,"tiempo_pasajero": 1,"tiempo_solicitud": 1,"tiempo_vacio": 1,"multiplicador_eod": 1,"accesibilidad": 1,"operador_mujer": 1,"inicio_eod": 1,"fin_eod": 1,"dist_pasajero_eod": 1,"tiempo_pasajero_eod": 1,"dist_solicitud": 1,"dist_vacío": 1,"dist_solicitud_eod": 1,"tiempo_solicitud_eod": 1,"dist_vacio_eod": 1,"tiempo_vacio_eod": 1}
```