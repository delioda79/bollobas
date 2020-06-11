# SEMOVI


### Kafka example messages

Topics:
- hecho_transito
```bash
$ kafka-console-producer --broker-list PLAINTEXT://localhost:9092 --topic semovi_beat_incidents\
> {"fecha_produccion": 1589473694068,"hecho_trans": null,"placa": null,"licencia": null,"distancia_viaje": null,"tiempo_viaje": null,"ubicacion": null,"tiempo_hecho": null}
```
- stats_operador
```bash
$ kafka-console-producer --broker-list PLAINTEXT://localhost:9092 --topic semovi_drivers_with_at_least_one_ride\
> {"fecha_produccion": 1589287539281,  "id_operador": "67151121-1022-4df0-abd3-09009eaa505f",  "genero": 1,  "cant_viajes": 1,  "tiempo_registro": 95,  "edad": "28-32",  "horas_conectado": "0-24",  "horas_viaje": "0-24",  "ingreso_totales": "$0-999"}
```
- viajes_agregados
```bash
$ kafka-console-producer --broker-list PLAINTEXT://localhost:9092 --topic semovi_beat_operation\
> {"fecha_produccion": 1589469747622,  "fecha": 1587106800000,  "id_proveedor": "BEAT",  "tot_viajes": 1,  "tot_veh_viaje": 1,  "tot_veh_disp": null,  "dist_pasajero": 6.87,  "tiempo_pasajero": 13,  "tiempo_solicitud": 13,  "dist_solicitud": 2.92,  "tiempo_vacio": null,  "multiplicador_eod": null,  "accesibilidad": null,  "operador_mujer": 0,  "inicio_eod": null,  "fin_eod": null,  "dist_pasajero_eod": null,  "tiempo_pasajero_eod": null,  "dist_vacio": null,  "dist_solicitud_eod": null,  "tiempo_solicitud_eod": null,  "dist_vacio_eod": null,  "tiempo_vacio_eod": null}
```