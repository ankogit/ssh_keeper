package styles

// Цвета приложения
const (
	// Основные цвета
	ColorPrimary    = "#7D56F4" // Фиолетовый - основной цвет
	ColorSecondary  = "#04B575" // Зеленый - вторичный цвет
	ColorWarning    = "#FFA500" // Оранжевый - предупреждение
	ColorError      = "#FF6B6B" // Красный - ошибка
	ColorSuccess    = "#51CF66" // Светло-зеленый - успех
	ColorText       = "#FFFFFF" // Белый - основной текст
	ColorMuted      = "#808080" // Серый - приглушенный текст
	ColorBackground = "#1A1A1A" // Темно-серый - фон (закомментирован)
	ColorContainer  = "#2D2D2D" // Серый - фон контейнеров (закомментирован)
)

// Стили границ
const (
	BorderRounded = "rounded" // Закругленные границы
)

// Стили выравнивания
const (
	AlignCenter = "center" // Центральное выравнивание
	AlignLeft   = "left"   // Левое выравнивание
)

// Стили текста
const (
	TextBold   = true // Жирный текст
	TextItalic = true // Курсив
)

// Размеры (минимальные значения)
const (
	MinWidth  = 80 // Минимальная ширина
	MinHeight = 24 // Минимальная высота
)

// Стили для различных компонентов
const (
	// Контейнер
	ContainerBorderWidth = 10 // Ширина рамки контейнера
	ContainerPadding     = 1  // Внутренние отступы контейнера

	// Заголовок
	HeaderPadding = 1 // Внутренние отступы заголовка
	HeaderMargin  = 1 // Внешние отступы заголовка

	// Контент
	ContentBorderWidth = 12 // Ширина рамки контента (2 символа с каждой стороны)
	ContentPadding     = 1  // Внутренние отступы контента

	// Список
	ListTitleMargin      = 1 // Внешние отступы заголовка списка
	ListPaginationMargin = 1 // Внешние отступы пагинации списка
)
