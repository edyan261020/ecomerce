package main

import (
	"database/sql"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var err error

func main() {
	// Conexión a la base de datos
	dsn := "root:@tcp(127.0.0.1:3306)/ecomerce"
	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Error al conectar a la base de datos:", err)
	}
	defer db.Close()

	// Verificar conexión
	err = db.Ping()
	if err != nil {
		log.Fatal("No se pudo conectar a la base de datos:", err)
	}

	r := gin.Default()
	r.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")
		c.Writer.Header().Set("Pragma", "no-cache")
		c.Writer.Header().Set("Expires", "0")
		c.Next()
	})

	// Servir archivos estáticos y templates
	r.Static("/static", "./static")
	r.LoadHTMLGlob("templates/*")

	// Páginas HTML
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "welcome.html", gin.H{"title": "Bienvenido a mi e-comerce"})
	})
	r.GET("/products", func(c *gin.Context) {
		c.HTML(http.StatusOK, "products.html", gin.H{"title": "Productos"})
	})
	r.GET("/login", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", gin.H{"title": "Iniciar Sesión"})
	})
	r.GET("/admin", func(c *gin.Context) {
		c.HTML(http.StatusOK, "admin.html", gin.H{"title": "Iniciar Sesión"})
	})

	r.GET("/users", func(c *gin.Context) {
		c.HTML(http.StatusOK, "users.html", gin.H{"title": "Iniciar Sesión"})
	})
	r.GET("/productsx", func(c *gin.Context) {
		c.HTML(http.StatusOK, "productsx.html", gin.H{"title": "Iniciar Sesión"})
	})
	r.GET("/sales", func(c *gin.Context) {
		c.HTML(http.StatusOK, "sales.html", gin.H{"title": "Iniciar Sesión"})
	})
	r.GET("/clients", func(c *gin.Context) {
		c.HTML(http.StatusOK, "clients.html", gin.H{"title": "Iniciar Sesión"})
	})

	// Rutas API REST
	api := r.Group("/api")
	{
		api.GET("/productos", getProductos)
		api.GET("/productos/:id", getProductoByID)
		api.POST("/productos", createProducto)
		api.PUT("/productos/:id", updateProducto)
		api.DELETE("/productos/:id", deleteProducto)

		api.POST("/login", login)

		api.GET("/usuarios", getUsuarios)
		api.POST("/usuarios", createUsuario)
		api.PUT("/usuarios/:cedula", updateUsuario)
		api.DELETE("/usuarios/:cedula", deleteUsuario)
		api.GET("/usuarios/:cedula", getUsuarioByCedulax)

		api.GET("/clientes", getClientes)
		api.POST("/clientes", createCliente)
		api.PUT("/clientes/:cedula", updateCliente)
		api.DELETE("/clientes/:cedula", deleteCliente)

		api.GET("/ventas", getVentas)
		api.GET("/ventas/max", getVentasMaxID)
		api.POST("/ventas", createVenta)
		api.PUT("/ventas/:id", updateVenta)
		api.DELETE("/ventas/:id", deleteVenta)

		api.GET("/detalles_ventas", getDetallesVenta)
		api.POST("/detalles_ventas", createDetalleVenta)
		api.GET("/detalles_ventas/:id_venta", getDetalleByVentaID)
		api.DELETE("/detalles_ventas/:id_venta/:id_producto", deleteDetalleVenta)

		api.GET("/ventasx", getVentasx)              // ahora soporta ?from=&to=
		api.GET("/ventas/resumen", getResumenVentas) // devuelve totalVentas & ingresos
		api.GET("/ventas/masvendido", getMasVendido) // top productos en rango
		// tus endpoints anteriores para POST/PUT/DELETE

	}

	// Iniciar servidor
	r.Run(":8080")
}

// Estructura de producto
type Producto struct {
	ID       int     `json:"id"`
	Name     string  `json:"name"`
	Category string  `json:"category"`
	Price    float64 `json:"price"`
	OldPrice float64 `json:"oldPrice"`
	Image    string  `json:"image"`
	Badge    string  `json:"badge"`
}
type Usuario struct {
	Cedula string `json:"cedula"`
	Nombre string `json:"nombre"`
	Rol    string `json:"rol"`
	User   string `json:"user"`
	Pass   string `json:"-"`
	Notas  string `json:"notas"`
}

// Cliente representa un registro en la tabla clientes
type Cliente struct {
	Cedula    string `json:"cedula"`
	Nombres   string `json:"nombres"`
	Apellidos string `json:"apellidos"`
	User      string `json:"user"`
}

type Venta struct {
	ID            int     `json:"id"`
	CedulaCliente string  `json:"cedula_cliente"`
	NombreCliente string  `json:"nombre_cliente"`
	Fecha         string  `json:"fecha"` // O puedes usar time.Time si manejas bien el parseo
	Total         float64 `json:"total"`
}
type DetalleVenta struct {
	IDVenta    int     `json:"id_venta"`
	IDProducto int     `json:"id_producto"`
	Detalle    string  `json:"detalle"`
	Cantidad   int     `json:"cantidad"`
	PrecioUni  float64 `json:"precio_uni"`
	SubTotal   float64 `json:"sub_total"`
}

// Obtener todos los productos
func getProductos(c *gin.Context) {
	rows, err := db.Query("SELECT * FROM productos")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var productos []Producto
	for rows.Next() {
		var p Producto
		if err := rows.Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.OldPrice, &p.Image, &p.Badge); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		productos = append(productos, p)
	}

	c.JSON(http.StatusOK, productos)
}

// ResumenVentas agrupa el total de ventas y sumatoria de ingresos
func getResumenVentas(c *gin.Context) {
	from, to := c.Query("from"), c.Query("to")
	var totalVentas int
	var ingresos float64

	err := db.QueryRow(
		`SELECT 
            COUNT(*) AS total_ventas, 
            IFNULL(SUM(total),0) AS ingresos 
         FROM ventas
         WHERE fecha BETWEEN ? AND ?`,
		from, to,
	).Scan(&totalVentas, &ingresos)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"totalVentas": totalVentas,
		"ingresos":    ingresos,
	})
}

// getVentas ahora filtra por rango opcional de fechas
func getVentasx(c *gin.Context) {
	from, to := c.Query("from"), c.Query("to")

	query := `SELECT id, cedula_cliente, nombre_cliente, fecha, total
              FROM ventas`
	args := []interface{}{}
	if from != "" && to != "" {
		query += " WHERE fecha BETWEEN ? AND ?"
		args = append(args, from, to)
	}
	query += " ORDER BY fecha"

	rows, err := db.Query(query, args...)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var ventas []Venta
	for rows.Next() {
		var v Venta
		if err := rows.Scan(&v.ID, &v.CedulaCliente, &v.NombreCliente, &v.Fecha, &v.Total); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ventas = append(ventas, v)
	}
	c.JSON(http.StatusOK, ventas)
}

// ProductoVend struct para el más vendido
type ProductoVend struct {
	IDProducto    int     `json:"id_producto"`
	Nombre        string  `json:"nombre"`
	TotalCantidad int     `json:"totalCantidad"`
	Ingresos      float64 `json:"ingresos"`
}

// getMasVendido devuelve los top 5 productos más vendidos en el rango
func getMasVendido(c *gin.Context) {
	from, to := c.Query("from"), c.Query("to")

	rows, err := db.Query(
		`SELECT 
            dv.id_producto,
            p.name,
            SUM(dv.cantidad) AS totalCantidad,
            SUM(dv.sub_total)   AS ingresos
         FROM detalle_venta dv
         JOIN productos p ON dv.id_producto = p.id
         WHERE dv.id_venta IN (
             SELECT id FROM ventas WHERE fecha BETWEEN ? AND ?
         )
         GROUP BY dv.id_producto, p.name
         ORDER BY totalCantidad DESC
         LIMIT 5`,
		from, to,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var stats []ProductoVend
	for rows.Next() {
		var pv ProductoVend
		if err := rows.Scan(&pv.IDProducto, &pv.Nombre, &pv.TotalCantidad, &pv.Ingresos); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		stats = append(stats, pv)
	}
	c.JSON(http.StatusOK, stats)
}

// GET /api/usuarios/:cedula
func getUsuarioByCedulax(c *gin.Context) {
	cedula := c.Param("cedula")

	var u Usuario
	err := db.QueryRow(
		"SELECT cedula, nombre, rol, user, notas FROM usuarios WHERE cedula = ?",
		cedula,
	).Scan(&u.Cedula, &u.Nombre, &u.Rol, &u.User, &u.Notas)

	if err == sql.ErrNoRows {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, u)
}

// Obtener un producto por ID
func getProductoByID(c *gin.Context) {
	id := c.Param("id")
	var p Producto
	err := db.QueryRow("SELECT * FROM productos WHERE id = ?", id).Scan(&p.ID, &p.Name, &p.Category, &p.Price, &p.OldPrice, &p.Image, &p.Badge)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Producto no encontrado"})
		return
	}
	c.JSON(http.StatusOK, p)
}

// Crear nuevo producto
func createProducto(c *gin.Context) {
	var p Producto
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stmt, err := db.Prepare("INSERT INTO productos(name, category, price, oldPrice, image, badge) VALUES (?, ?, ?, ?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	res, err := stmt.Exec(p.Name, p.Category, p.Price, p.OldPrice, p.Image, p.Badge)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	id, _ := res.LastInsertId()
	p.ID = int(id)
	c.JSON(http.StatusCreated, p)
}

// Actualizar un producto
func updateProducto(c *gin.Context) {
	id := c.Param("id")
	var p Producto
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	stmt, err := db.Prepare("UPDATE productos SET name=?, category=?, price=?, oldPrice=?, image=?, badge=? WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = stmt.Exec(p.Name, p.Category, p.Price, p.OldPrice, p.Image, p.Badge, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	p.ID = atoi(id)
	c.JSON(http.StatusOK, p)
}

// Eliminar un producto
func deleteProducto(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM productos WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Producto eliminado correctamente"})
}

// POST /api/login
func login(c *gin.Context) {
	var creds struct {
		User string `json:"user"` // para admin: user; para cliente: cedula
		Pass string `json:"pass"` // para admin: pass; para cliente: user (contraseña)
	}
	if err := c.ShouldBindJSON(&creds); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	type auth struct {
		StoredPass string
		Role       string
		Cedula     string
		Nombres    string
	}
	var a auth
	var err error

	// 1) Intentar en tabla usuarios (admin)
	err = db.QueryRow(
		"SELECT pass, 'admin' AS role, '' AS cedula, '' AS nombres FROM usuarios WHERE user = ?",
		creds.User,
	).Scan(&a.StoredPass, &a.Role, &a.Cedula, &a.Nombres)

	if err == sql.ErrNoRows {
		// 2) Intentar en tabla clientes
		err = db.QueryRow(
			"SELECT user, 'client' AS role, cedula, nombres FROM clientes WHERE cedula = ?",
			creds.User,
		).Scan(&a.StoredPass, &a.Role, &a.Cedula, &a.Nombres)
		if err == sql.ErrNoRows {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña inválidos"})
			return
		} else if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 3) Comparar contraseñas
	if a.StoredPass != creds.Pass {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Usuario o contraseña inválidos"})
		return
	}

	// 4) Éxito: devolvemos role y, si es client, cedula/nombres
	c.JSON(http.StatusOK, gin.H{
		"user":    creds.User,
		"role":    a.Role,
		"cedula":  a.Cedula,
		"nombres": a.Nombres,
	})
}

// Obtener todos los usuarios
func getUsuarios(c *gin.Context) {
	rows, err := db.Query("SELECT cedula, nombre, rol, user, notas FROM usuarios")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var usuarios []Usuario
	for rows.Next() {
		var u Usuario
		if err := rows.Scan(&u.Cedula, &u.Nombre, &u.Rol, &u.User, &u.Notas); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		usuarios = append(usuarios, u)
	}

	c.JSON(http.StatusOK, usuarios)
}

// Obtener un usuario por cédula
func getUsuarioByCedula(c *gin.Context) {
	cedula := c.Param("cedula")
	var u Usuario
	err := db.QueryRow(
		"SELECT cedula, nombre, rol, user, notas FROM usuarios WHERE cedula = ?",
		cedula,
	).Scan(&u.Cedula, &u.Nombre, &u.Rol, &u.User, &u.Notas)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Usuario no encontrado"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// Crear nuevo usuario
func createUsuario(c *gin.Context) {
	var u Usuario
	type payload struct {
		Cedula string `json:"cedula"`
		Nombre string `json:"nombre"`
		Rol    string `json:"rol"`
		User   string `json:"user"`
		Pass   string `json:"pass"`
		Notas  string `json:"notas"`
	}
	var in payload

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare(
		"INSERT INTO usuarios(cedula, nombre, rol, user, pass, notas) VALUES (?, ?, ?, ?, ?, ?)",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(in.Cedula, in.Nombre, in.Rol, in.User, in.Pass, in.Notas)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	u = Usuario{
		Cedula: in.Cedula,
		Nombre: in.Nombre,
		Rol:    in.Rol,
		User:   in.User,
		Notas:  in.Notas,
	}
	c.JSON(http.StatusCreated, u)
}

// Actualizar un usuario
func updateUsuario(c *gin.Context) {
	cedula := c.Param("cedula")
	type payload struct {
		Nombre string `json:"nombre"`
		Rol    string `json:"rol"`
		User   string `json:"user"`
		Pass   string `json:"pass"`
		Notas  string `json:"notas"`
	}
	var in payload

	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare(
		"UPDATE usuarios SET nombre=?, rol=?, user=?, pass=?, notas=? WHERE cedula=?",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(in.Nombre, in.Rol, in.User, in.Pass, in.Notas, cedula)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cedula": cedula,
		"nombre": in.Nombre,
		"rol":    in.Rol,
		"user":   in.User,
		"notas":  in.Notas,
	})
}

// Eliminar un usuario
func deleteUsuario(c *gin.Context) {
	cedula := c.Param("cedula")
	stmt, err := db.Prepare("DELETE FROM usuarios WHERE cedula = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cedula)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Usuario eliminado correctamente"})
}

// Obtener todos los clientes
func getClientes(c *gin.Context) {
	rows, err := db.Query("SELECT cedula, nombres, apellidos, user FROM clientes")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var clientes []Cliente
	for rows.Next() {
		var cl Cliente
		if err := rows.Scan(&cl.Cedula, &cl.Nombres, &cl.Apellidos, &cl.User); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		clientes = append(clientes, cl)
	}

	c.JSON(http.StatusOK, clientes)
}

// Obtener un cliente por cédula
func getClienteByCedula(c *gin.Context) {
	cedula := c.Param("cedula")
	var cl Cliente
	err := db.QueryRow(
		"SELECT cedula, nombres, apellidos, user FROM clientes WHERE cedula = ?",
		cedula,
	).Scan(&cl.Cedula, &cl.Nombres, &cl.Apellidos, &cl.User)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Cliente no encontrado"})
		return
	}
	c.JSON(http.StatusOK, cl)
}

// Crear nuevo cliente
func createCliente(c *gin.Context) {
	type payload struct {
		Cedula    string `json:"cedula"`
		Nombres   string `json:"nombres"`
		Apellidos string `json:"apellidos"`
		User      string `json:"user"`
		Pass      string `json:"pass"`
	}
	var in payload
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare(
		"INSERT INTO clientes(cedula, nombres, apellidos, user, pass) VALUES (?, ?, ?, ?, ?)",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(in.Cedula, in.Nombres, in.Apellidos, in.User, in.Pass)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	cl := Cliente{
		Cedula:    in.Cedula,
		Nombres:   in.Nombres,
		Apellidos: in.Apellidos,
		User:      in.User,
	}
	c.JSON(http.StatusCreated, cl)
}

// Actualizar un cliente
func updateCliente(c *gin.Context) {
	cedula := c.Param("cedula")
	type payload struct {
		Nombres   string `json:"nombres"`
		Apellidos string `json:"apellidos"`
		User      string `json:"user"`
		Pass      string `json:"pass"`
	}
	var in payload
	if err := c.ShouldBindJSON(&in); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare(
		"UPDATE clientes SET nombres=?, apellidos=?, user=?, pass=? WHERE cedula=?",
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(in.Nombres, in.Apellidos, in.User, in.Pass, cedula)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"cedula":    cedula,
		"nombres":   in.Nombres,
		"apellidos": in.Apellidos,
		"user":      in.User,
	})
}

// Eliminar un cliente
func deleteCliente(c *gin.Context) {
	cedula := c.Param("cedula")
	stmt, err := db.Prepare("DELETE FROM clientes WHERE cedula = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(cedula)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cliente eliminado correctamente"})
}

func getVentas(c *gin.Context) {
	rows, err := db.Query("SELECT id, cedula_cliente, nombre_cliente, fecha, total FROM ventas")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var ventas []Venta
	for rows.Next() {
		var v Venta
		if err := rows.Scan(&v.ID, &v.CedulaCliente, &v.NombreCliente, &v.Fecha, &v.Total); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		ventas = append(ventas, v)
	}

	c.JSON(http.StatusOK, ventas)
}

// GET /api/ventas/max
func getVentasMaxID(c *gin.Context) {
	// Ejecuta la consulta para obtener el máximo id
	row := db.QueryRow("SELECT MAX(id) AS ultimo_id FROM ventas")

	var ultimoID int
	if err := row.Scan(&ultimoID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Devuelve solo el campo ultimo_id
	c.JSON(http.StatusOK, gin.H{"id": ultimoID})
}

func getVentaByID(c *gin.Context) {
	id := c.Param("id")
	var v Venta
	err := db.QueryRow("SELECT id, cedula_cliente, nombre_cliente, fecha, total FROM ventas WHERE id = ?", id).
		Scan(&v.ID, &v.CedulaCliente, &v.NombreCliente, &v.Fecha, &v.Total)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Venta no encontrada"})
		return
	}

	c.JSON(http.StatusOK, v)
}
func createVenta(c *gin.Context) {
	var v Venta
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("INSERT INTO ventas(cedula_cliente, nombre_cliente, total) VALUES (?, ?, ?)")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	res, err := stmt.Exec(v.CedulaCliente, v.NombreCliente, v.Total)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	id, _ := res.LastInsertId()
	v.ID = int(id)

	// Puedes recuperar la fecha generada automáticamente, si deseas
	err = db.QueryRow("SELECT fecha FROM ventas WHERE id = ?", id).Scan(&v.Fecha)
	if err != nil {
		v.Fecha = "fecha no disponible"
	}

	c.JSON(http.StatusCreated, v)
}
func updateVenta(c *gin.Context) {
	id := c.Param("id")
	var v Venta
	if err := c.ShouldBindJSON(&v); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare("UPDATE ventas SET cedula_cliente=?, nombre_cliente=?, total=? WHERE id=?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(v.CedulaCliente, v.NombreCliente, v.Total, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	v.ID = atoi(id)
	c.JSON(http.StatusOK, v)
}

func deleteVenta(c *gin.Context) {
	id := c.Param("id")
	stmt, err := db.Prepare("DELETE FROM ventas WHERE id = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Venta eliminada correctamente"})
}

func getDetallesVenta(c *gin.Context) {
	rows, err := db.Query("SELECT id_venta, id_producto, detalle, cantidad, precio_uni, sub_total FROM detalle_venta")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var detalles []DetalleVenta
	for rows.Next() {
		var d DetalleVenta
		if err := rows.Scan(&d.IDVenta, &d.IDProducto, &d.Detalle, &d.Cantidad, &d.PrecioUni, &d.SubTotal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		detalles = append(detalles, d)
	}

	c.JSON(http.StatusOK, detalles)
}

func getDetalleByVentaID(c *gin.Context) {
	id := c.Param("id_venta")
	rows, err := db.Query("SELECT id_venta, id_producto, detalle, cantidad, precio_uni, sub_total FROM detalle_venta WHERE id_venta = ?", id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var detalles []DetalleVenta
	for rows.Next() {
		var d DetalleVenta
		if err := rows.Scan(&d.IDVenta, &d.IDProducto, &d.Detalle, &d.Cantidad, &d.PrecioUni, &d.SubTotal); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		detalles = append(detalles, d)
	}

	if len(detalles) == 0 {
		c.JSON(http.StatusNotFound, gin.H{"message": "No se encontraron detalles para esta venta"})
		return
	}

	c.JSON(http.StatusOK, detalles)
}
func createDetalleVenta(c *gin.Context) {
	var d DetalleVenta
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare(`INSERT INTO detalle_venta(id_venta, id_producto, detalle, cantidad, precio_uni, sub_total)
        VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(d.IDVenta, d.IDProducto, d.Detalle, d.Cantidad, d.PrecioUni, d.SubTotal)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, d)
}
func updateDetalleVenta(c *gin.Context) {
	idVenta := c.Param("id_venta")
	idProducto := c.Param("id_producto")

	var d DetalleVenta
	if err := c.ShouldBindJSON(&d); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	stmt, err := db.Prepare(`UPDATE detalle_venta SET detalle=?, cantidad=?, precio_uni=?, sub_total=?
        WHERE id_venta=? AND id_producto=?`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(d.Detalle, d.Cantidad, d.PrecioUni, d.SubTotal, idVenta, idProducto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	d.IDVenta = atoi(idVenta)
	d.IDProducto = atoi(idProducto)
	c.JSON(http.StatusOK, d)
}
func deleteDetalleVenta(c *gin.Context) {
	idVenta := c.Param("id_venta")
	idProducto := c.Param("id_producto")

	stmt, err := db.Prepare("DELETE FROM detalle_venta WHERE id_venta = ? AND id_producto = ?")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer stmt.Close()

	_, err = stmt.Exec(idVenta, idProducto)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Detalle de venta eliminado correctamente"})
}

// Convertir string a int
func atoi(s string) int {
	i, _ := strconv.Atoi(s)
	return i
}
