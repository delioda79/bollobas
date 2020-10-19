# SEMOVI Technical Documentation
 

## Version: 1.0.0

**Contact information:**  Beat  
**License:** BEAT Mobility Services


# /semovi/1.0.0/hecho_transito

### Method: GET
### Summary:

This method returns all the traffic incidents data. The Response Data will be available for consultation during the ten business days of the month. This means that e.g. from April 01 through April 10, aggregated data only for March should be available to fetch. It is important to mention that after this period of time no data will be exposed.
The API have a method to request data by date ranges using the parameters from, to.
The response is a paginated response. By default you will receive 10 items per page. 

### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| from | query | start date (YmdTHis) (ex. 20200504T083521) | No | integer |
| to | query | end date (YmdTHis) (ex. 20200530T150020) | No | integer |
| page | query | Page Value | No | integer |
| size | query | Size Value | No | integer |

### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [trafficIncident](#trafficincident) |
| 400 | Bad Request | [error](#error) |
| 500 | Internal Server Error | [error](#error) |

# /semovi/1.0.0/stats_operador

This method returns all the aggregated operator stats data. The Response Data will be available for consultation during the ten business days of the month. This means that e.g. from April 01 through April 10, aggregated data only for March should be available to fetch. It is important to mention that after this period of time no data should be exposed.
The API have a method to request data by date ranges using the parameters from, to.
The response is a paginated response. By default you will receive 10 items. 

### Method: GET
##### Summary:

Get return all the operator stats

##### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| from | query | start date (YmdTHis) (ex. 20200504T083521) | No | integer |
| to | query | end date (YmdTHis) (ex. 20200530T150020) | No | integer |
| page | query | Page Value | No | integer |
| size | query | Size Value | No | integer |

##### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [trafficIncident](#trafficincident) |
| 400 | Bad Request | [error](#error) |
| 500 | Internal Server Error | [error](#error) |

# /semovi/1.0.0/viajes_agregados

### Method: GET
### Summary:

This method returns all the the aggregated trips data. The Response Data will be available for consultation during the ten business days of the month. This means that e.g. from April 01 through April 10, aggregated data only for March should be available to fetch. It is important to mention that after this period of time no data should be exposed.
The API have a method to request data by date ranges using the parameters from, to.
The response is a paginated response. By default you will receive 10 items. 

### Parameters

| Name | Located in | Description | Required | Schema |
| ---- | ---------- | ----------- | -------- | ---- |
| from | query | start date (YmdTHis) (ex. 20200504T083521) | No | integer |
| to | query | end date (YmdTHis) (ex. 20200530T150020) | No | integer |
| page | query | Page Value | No | integer |
| size | query | Size Value | No | integer |

### Responses

| Code | Description | Schema |
| ---- | ----------- | ------ |
| 200 | OK | [trafficIncident](#trafficincident) |
| 400 | Bad Request | [error](#error) |
| 500 | Internal Server Error | [error](#error) |

### Models


#### AggregatedTrips

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| accesibilidad | number |  | No |
| dist_pasajero | number |  | No |
| dist_pasajero_eod | number |  | No |
| dist_solicitud | number |  | No |
| dist_solicitud_eod | number |  | No |
| dist_vacio_eod | number |  | No |
| dist_vacío | number |  | No |
| fecha | string |  | No |
| fin_eod | integer |  | No |
| id | integer |  | No |
| id_proveedor | string |  | No |
| inicio_eod | integer |  | No |
| multiplicador_eod | number |  | No |
| operador_mujer | number |  | No |
| tiempo_pasajero | number |  | No |
| tiempo_pasajero_eod | integer |  | No |
| tiempo_solicitud | number |  | No |
| tiempo_solicitud_eod | number |  | No |
| tiempo_vacio | number |  | No |
| tiempo_vacio_eod | number |  | No |
| tot_veh_disp | integer |  | No |
| tot_veh_viaje | integer |  | No |
| tot_viajes | integer |  | No |

#### Error

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| error | string |  | No |

#### OperatorStats

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| cant_viajes | integer |  | No |
| edad | string |  | No |
| fecha_produccion | string |  | No |
| genero | integer |  | No |
| horas_conectado | string |  | No |
| horas_viaje | string |  | No |
| id | integer |  | No |
| id_operador | string |  | No |
| ingreso_totales | string |  | No |
| tiempo_registro | integer |  | No |

#### TrafficIncident

| Name | Type | Description | Required |
| ---- | ---- | ----------- | -------- |
| distancia_viaje | string |  | No |
| hecho_trans | integer |  | No |
| id | integer |  | No |
| licencia | string |  | No |
| placa | string |  | No |
| tiempo_hecho | string |  | No |
| tiempo_viaje | string |  | No |
| ubicación | string |  | No |


# Examples:

### Testing Environment
```
curl --location --request GET 'https://qamx-rest.k8s.sandbox.thebeat.co/semovi/1.0.0/viajes_agregados?from=20200504T083521&to=20200530T150020&page=1&size=10' \
--header 'X-SECURE-ROUTE-KEY: zp4V8pFB2JmmjDA2Wk8Pq1G90l'

curl --location --request GET 'https://qamx-rest.k8s.sandbox.thebeat.co/semovi/1.0.0/stats_operador?from=20200504T083521&to=20200530T150020&page=1&size=10' \
--header 'X-SECURE-ROUTE-KEY: zp4V8pFB2JmmjDA2Wk8Pq1G90l'

curl --location --request GET 'https://qamx-rest.k8s.sandbox.thebeat.co/semovi/1.0.0/hecho_transito?from=20200504T083521&to=20200530T150020&page=1&size=10' \
--header 'X-SECURE-ROUTE-KEY: zp4V8pFB2JmmjDA2Wk8Pq1G90l'
```

### Live Environment
```
- Production:
curl --location --request GET 'https://rest.mexico.thebeat.co/semovi/1.0.0/viajes_agregados?from=20200504T083521&to=20200530T150020&page=1&size=10' \
--header 'X-SECURE-ROUTE-KEY: zp4V8pFB2JmmjDA2Wk8Pq1G90l'

curl --location --request GET 'https://rest.mexico.thebeat.co/semovi/1.0.0/stats_operador?from=20200504T083521&to=20200530T150020&page=1&size=10' \
--header 'X-SECURE-ROUTE-KEY: zp4V8pFB2JmmjDA2Wk8Pq1G90l'

curl --location --request GET 'https://rest.mexico.thebeat.co/semovi/1.0.0/hecho_transito?from=20200504T083521&to=20200530T150020&page=1&size=10' \
--header 'X-SECURE-ROUTE-KEY: zp4V8pFB2JmmjDA2Wk8Pq1G90l'
```

Respnse

```
{
    "Data": [
        {
            "id": 1,
            "tiempo_hecho": "2020-05-12T10:27:37",
            "hecho_trans": 2,
            "placa": "ABC-123",
            "licencia": "C12345678",
            "distancia_viaje": "15-20",
            "tiempo_viaje": "100-200",
            "ubicación": "38.0088261,23.8042912"
        },
        {
            "id": 341,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "DHS 5956",
            "licencia": "CHYPKEKA27167",
            "distancia_viaje": "36000-38999",
            "tiempo_viaje": "30-34",
            "ubicación": "19.43436691571364,-99.1337055637379"
        },
        {
            "id": 342,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "TGN 4504",
            "licencia": "OYNKMSQV88979",
            "distancia_viaje": "15000-17999",
            "tiempo_viaje": "5-9",
            "ubicación": "19.432425565935215,-99.13466988859915"
        },
        {
            "id": 343,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "KIV 4374",
            "licencia": "DOQBHJQK90378",
            "distancia_viaje": "24000-26999",
            "tiempo_viaje": "75-79",
            "ubicación": "19.431146138522617,-99.13361748517156"
        },
        {
            "id": 344,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "PNK 9509",
            "licencia": "JQOXKQUN49734",
            "distancia_viaje": "45000-47999",
            "tiempo_viaje": "5-9",
            "ubicación": "19.431654302566397,-99.13216951440148"
        },
        {
            "id": 345,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "STB 3407",
            "licencia": "ULHGHNOS45710",
            "distancia_viaje": "51000_o_mas",
            "tiempo_viaje": "65-69",
            "ubicación": "19.432497226467394,-99.13183816200464"
        },
        {
            "id": 346,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "AFK 4272",
            "licencia": "XCRXLPLN76466",
            "distancia_viaje": "51000_o_mas",
            "tiempo_viaje": "90-94",
            "ubicación": "19.434594922539087,-99.13557540795294"
        },
        {
            "id": 347,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "PLU 2745",
            "licencia": "DCQTGMLZ03372",
            "distancia_viaje": "33000-35999",
            "tiempo_viaje": "70-74",
            "ubicación": "19.43162399115729,-99.13152412635662"
        },
        {
            "id": 348,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "IGK 1613",
            "licencia": "TTWLHSYL82238",
            "distancia_viaje": "39000-41999",
            "tiempo_viaje": "60-64",
            "ubicación": "19.42931933555282,-99.13452386061329"
        },
        {
            "id": 349,
            "tiempo_hecho": "2020-05-05T08:40:52",
            "hecho_trans": 1,
            "placa": "RCS 3572",
            "licencia": "OGBREWFK75147",
            "distancia_viaje": "0-2999",
            "tiempo_viaje": "75-79",
            "ubicación": "19.431592927126673,-99.13286638492379"
        }
    ],
    "Meta": {
        "TotalCount": 262,
        "TotalPages": 27,
        "PageSize": 10,
        "CurrentPage": 1
    }
}
```