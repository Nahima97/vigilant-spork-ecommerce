# Vigilant-Spork — E-Commerce API (Golang)

### 📌 Project Title & Overview

Vigilant-Spork is a production-ready E-Commerce REST API built with Golang, PostgreSQL, GORM, and JWT Authentication.

It supports core e-commerce functionality including:
- User registration & login
- Role-based access (admin vs. customer)
- Product management
- Product reviews
- Full shopping cart CRUD functionality

This API exists to provide a clean, modular, and scalable backend template suitable for learning, portfolio projects, and real-world development.

### 📝 Description / Purpose

E-commerce platforms require secure authentication, structured product management, and a reliable shopping cart system.

This project solves that by providing:
- A clean architecture that separates concerns (Handlers, Services, Repositories)
- Safe, secure authentication with JWT
- Cart calculations using integer currency math to avoid floating-point errors
- Reusable API patterns suitable for any modern backend project

It is designed for:
- Students learning backend development
- Engineers building a scalable e-commerce backend
- Teams wanting an example of clean Golang architecture
- Portfolios needing a production-style API project

## ⚡ Features

### 🔐 Authentication

- Register user (role restricted to admin or customer)
- Login with JWT token
- Protected routes using middleware

### 🛍️ Products

- Create, update, delete products (admin-only)
- List products with:
- Pagination
- Category filters
- Price filters
- View product details with:
- Average rating
- Reviews
- “No user reviews” fallback message

### ⭐ Reviews

- Customers can leave reviews
- Rating automatically updates per product

### 🛒 Shopping Cart

Complete CRUD support:
- Add To Cart
- View Cart (dynamically calculates total using current prices)
- Update Item Quantity
- Remove Item

## 🛠 Tech Stack / Requirements

### 🧰 Languages & Libraries

- Golang 1.20+
- Gorilla Mux for routing
- GORM for ORM
- UUID package (gofrs/uuid)
- JWT-Go for authentication

## 🏆 Credits / Acknowledgments

### Special thanks to:

- Golang & GORM Communities for documentation and examples
- PostgreSQL for robust and open-source database technology
- Project authors Nahima Sultana, Safa Aydarus, and Sevi Hazim for building a clean, production-level backend
