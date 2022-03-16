#!/usr/bin/bash

sleep 1

curl -H "Content-Type: application/json" -d '{"match": "FirstBrand","reward": 10, "reward_type": "%"}' -X POST http://localhost:8888/api/goods
curl -H "Content-Type: application/json" -d '{"order": "1149", "goods": [{"description": "Чайник FirstBrand", "price": 7000}, {"description": "Ноутбук FirstBrand", "price": 35000}]}' -X POST http://localhost:8888/api/orders
curl -X GET http://localhost:8888/api/orders/1149

curl -H "Content-Type: application/json" -d '{"match": "SecondBrand","reward": 20, "reward_type": "%"}' -X POST http://localhost:8888/api/goods
curl -H "Content-Type: application/json" -d '{"order": "2253", "goods": [{"description": "Телефон SecondBrand", "price": 17000}, {"description": "Домашний кинотеатр SecondBrand", "price": 52000}]}' -X POST http://localhost:8888/api/orders
curl -X GET http://localhost:8888/api/orders/2253

curl -H "Content-Type: application/json" -d '{"match": "ThirdBrand","reward": 5, "reward_type": "%"}' -X POST http://localhost:8888/api/goods
curl -H "Content-Type: application/json" -d '{"order": "3376", "goods": [{"description": "Фонарь ThirdBrand", "price": 2000}, {"description": "Планшет ThirdBrand", "price": 5000}]}' -X POST http://localhost:8888/api/orders
curl -X GET http://localhost:8888/api/orders/3376