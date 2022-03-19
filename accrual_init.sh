#!/usr/bin/bash

sleep 1

curl -H "Content-Type: application/json" -d '{"match": "FirstBrand","reward": 10, "reward_type": "%"}' -X POST http://localhost:8888/api/goods
curl -H "Content-Type: application/json" -d '{"order": "1149", "goods": [{"description": "Чайник FirstBrand", "price": 700}, {"description": "Ноутбук FirstBrand", "price": 3500}]}' -X POST http://localhost:8888/api/orders
curl -X GET http://localhost:8888/api/orders/1149

curl -H "Content-Type: application/json" -d '{"match": "SecondBrand","reward": 20, "reward_type": "%"}' -X POST http://localhost:8888/api/goods
curl -H "Content-Type: application/json" -d '{"order": "2253", "goods": [{"description": "Телефон SecondBrand", "price": 1700}, {"description": "Домашний кинотеатр SecondBrand", "price": 5200}]}' -X POST http://localhost:8888/api/orders
curl -X GET http://localhost:8888/api/orders/2253

curl -H "Content-Type: application/json" -d '{"match": "ThirdBrand","reward": 5, "reward_type": "%"}' -X POST http://localhost:8888/api/goods
curl -H "Content-Type: application/json" -d '{"order": "3376", "goods": [{"description": "Фонарь ThirdBrand", "price": 200}, {"description": "Планшет ThirdBrand", "price": 500}]}' -X POST http://localhost:8888/api/orders
curl -X GET http://localhost:8888/api/orders/3376

curl -H "Content-Type: application/json" -d '{"order": "5587", "goods": [{"description": "Пароварка FifthBrand", "price": 2300}, {"description": "Клавиатура FourthBrand", "price": 150}]}' -X POST http://localhost:8888/api/orders
curl -X GET http://localhost:8888/api/orders/5587
