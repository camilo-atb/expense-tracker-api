
# 📝 Expense Tracker API (RESTful)

Expense Tracker API es una **API RESTful desarrollada en Go** para el registro y gestión de movimientos financieros personales.
Este proyecto fue creado con fines **educativos**, enfocado en **arquitectura backend**, **separación de responsabilidades** y **buenas prácticas en el diseño de APIs**, simulando un entorno profesional real.

La idea base del proyecto fue tomada de:
https://roadmap.sh/backend/project-ideas  
Sin embargo, **no fue desarrollado siguiendo estrictamente las indicaciones de roadmap.sh**, sino como un ejercicio propio de diseño y aprendizaje.

---

## 🎯 Objetivos del proyecto

- Practicar lógica de programación en Go
- Aplicar una arquitectura escalable y mantenible
- Implementar una API REST siguiendo buenas prácticas
- Separar claramente las capas: **handler**, **service** y **repository**

---

## ✨ Funcionalidades

- Registrar movimientos financieros
- Editar movimientos
- Eliminar movimientos (hard delete)
- Eliminar movimientos de forma lógica (soft delete)
- Listar todos los movimientos
- Filtrar movimientos por tipo
- Filtrar movimientos por rango de fechas
- Obtener totales por tipo y rango de fechas
- Calcular ingresos netos

Cada movimiento contiene:

- ID  
- Description  
- Amount  
- Category  
- Type  
- Date  
- Status  

---

## 🧑‍💻 Endpoints

### Obtener todos los movimientos
```http
GET /api/transactions
```

Ejemplo:
```http
http://localhost:8080/api/transactions
```

---

### Registrar un movimiento
```http
POST /api/transactions
```

Body:
```json
{
  "description": "monto 13",
  "amount": 50000,
  "category": "ocio",
  "type": "income"
}
```

---

### Editar un movimiento
```http
PATCH /api/transactions/{id}
```

Ejemplo:
```http
http://localhost:8080/api/transactions/6
```

---

### Eliminar un movimiento (hard delete)
```http
DELETE /api/transactions/{id}?mode=hard
```

---

### Obtener movimientos por rango de fechas
```http
GET /api/transactions?from=2026-01-01&to=2026-01-10
```

---

### Obtener movimientos por tipo y fechas
```http
GET /api/transactions?type=income&from=2026-01-01&to=2026-01-10
```

---

### Obtener total por tipo y fechas
```http
GET /api/summary/type?type=expense&from=2026-01-01&to=2026-01-30
```

---

### Obtener ingresos netos
```http
GET /api/summary/net
```

---

### Obtener movimiento por ID
```http
GET /api/transactions/{id}
```

---

## 🗂️ Estructura del proyecto

```
.
├── cmd/
│   └── api/
│       └── main.go
│
├── internal/
│   ├── expense/
│   │   ├── handler.go       // Capa HTTP (request/response)
│   │   ├── service.go       // Lógica de negocio
│   │   ├── repository.go    // Acceso a datos
│   │   ├── model.go         // Entidades del dominio
│   │   ├── create_dto.go    // DTO de creación
│   │   └── update_dto.go    // DTO de actualización
│   │
│   └── shared/
│       ├── httpx/           // Helpers HTTP (responses, errors)
│       └── database/        // Conexión a base de datos
│
├── go.mod
├── .env.example
├── README.md
├── .gitignore
└── LICENSE
```

---

## 📦 Dependencias principales

Dependencias utilizadas en el proyecto:

- **chi** – Router HTTP ligero y idiomático
- **pgx** – Driver PostgreSQL de alto rendimiento
- **godotenv** – Carga de variables de entorno

Las dependencias indirectas son gestionadas automáticamente por **Go Modules**.

---

## 🚀 Ejecución del proyecto

1. Clonar el repositorio
2. Configurar el archivo `.env` basado en `.env.example`
3. Ejecutar el servidor:

```bash
go run cmd/api/main.go
```

Servidor disponible en:
```http
http://localhost:8080
```

---

## 📌 Notas finales

- Proyecto con fines educativos
- Enfocado en buenas prácticas de backend
- Base sólida para futuras mejoras

Este proyecto forma parte de una serie de **proyectos de nivel inicial (nivel 1)**.
En niveles posteriores se planea agregar:

- Autenticación y autorización
- Manejo de usuarios
- Seguridad
- Validaciones avanzadas
- Tests automatizados

✍️ Desarrollado en Go como ejercicio práctico de aprendizaje.
