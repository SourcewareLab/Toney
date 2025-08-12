package config

type StylesConfig struct {
	Text             string           `mapstructure:"text"`
	Background       string           `mapstructure:"background"`
	Border           string           `mapstructure:"border"`
	FocusedBorder    string           `mapstructure:"focused_border"`
	MenuSelectedBg   string           `mapstructure:"menu_selected_bg"`
	MenuSelectedText string           `mapstructure:"menu_selected_text"`
	ErrorBg          string           `mapstructure:"error_background"` // TODO: Docs
	ErrorText        string           `mapstructure:"error_text"`       // TODO: Docs
	Icons            IconsConfig      `mapstructure:"icons"`
	Renderer         RendererConfig   `mapstructure:"renderer"`
	TaskStyles       TaskStylesConfig `mapstructure:"task_styles"`
}

type IconsConfig struct {
	FolderIcon string    `mapstructure:"folder_icon"`
	FileIcon   string    `mapstructure:"file_icon"`
	TaskIcons  TaskIcons `mapstructure:"task_icons"`
}

type TaskIcons struct {
	CompletedIcon string `mapstructure:"completed_icon"`
	PendingIcon   string `mapstructure:"pending_icon"`
	AbandonedIcon string `mapstructure:"abandoned_icon"`
	StartedIcon   string `mapstructure:"started_icon"`
}

type RendererConfig struct {
	ItemPrefix        string         `mapstructure:"item_prefix"`
	EnumerationPrefix string         `mapstructure:"enumeration_prefix"`
	Document          BlockStyle     `mapstructure:"document"`
	BlockQuote        BlockStyle     `mapstructure:"blockquote"`
	Heading           HeadingStyle   `mapstructure:"heading"`
	HorizontalRule    InlineStyle    `mapstructure:"horizontal_rule"`
	Paragraph         BlockStyle     `mapstructure:"paragraph"`
	List              ListStyle      `mapstructure:"list"`
	Enumeration       InlineStyle    `mapstructure:"enumeration"`
	Link              InlineStyle    `mapstructure:"link"`
	Image             InlineStyle    `mapstructure:"image"`
	Emph              InlineStyle    `mapstructure:"emph"`
	Strong            InlineStyle    `mapstructure:"strong"`
	Strikethrough     InlineStyle    `mapstructure:"strikethrough"`
	Code              BlockStyle     `mapstructure:"code"`
	CodeBlock         CodeBlockStyle `mapstructure:"code_block"`
	Table             TableStyle     `mapstructure:"table"`
}

type BlockStyle struct {
	Format      string `mapstructure:"format"`
	BlockPrefix string `mapstructure:"block_prefix"`
	BlockSuffix string `mapstructure:"block_suffix"`
	Prefix      string `mapstructure:"prefix"`
	Suffix      string `mapstructure:"suffix"`
	IndentToken string `mapstructure:"indent_token"`
	Margin      int    `mapstructure:"margin"`
	Padding     int    `mapstructure:"padding"`
	Indent      int    `mapstructure:"indent"`
	Color       string `mapstructure:"color"`
	Background  string `mapstructure:"background"`
	Bold        bool   `mapstructure:"bold"`
	Italic      bool   `mapstructure:"italic"`
	Underline   bool   `mapstructure:"underline"`
}

type HeadingStyle struct {
	Base   BlockStyle   `mapstructure:"base"`
	Levels []BlockStyle `mapstructure:"levels"`
}

type ListStyle struct {
	LevelIndent int         `mapstructure:"indent"`
	Styles      InlineStyle `mapstructure:"styles"`
	Task        *TaskStyle  `mapstructure:"task,omitempty"`
}

type TaskStyle struct {
	Ticked   string `mapstructure:"ticked"`
	Unticked string `mapstructure:"unticked"`
}

type InlineStyle struct {
	BlockPrefix   string `mapstructure:"block_prefix"`
	BlockSuffix   string `mapstructure:"block_suffix"`
	Prefix        string `mapstructure:"prefix"`
	Suffix        string `mapstructure:"suffix"`
	Color         string `mapstructure:"color"`
	Background    string `mapstructure:"background"`
	Bold          bool   `mapstructure:"bold"`
	Italic        bool   `mapstructure:"italic"`
	Underline     bool   `mapstructure:"underline"`
	Strikethrough bool   `mapstructure:"strikethrough"`
	Format        string `mapstructure:"format"`
}

type CodeBlockStyle struct {
	BlockStyle          `mapstructure:",squash"`
	Theme               string      `mapstructure:"theme"`
	Text                InlineStyle `mapstructure:"text"`
	Error               InlineStyle `mapstructure:"error"`
	Comment             InlineStyle `mapstructure:"comment"`
	CommentPreproc      InlineStyle `mapstructure:"comment_preproc"`
	Keyword             InlineStyle `mapstructure:"keyword"`
	KeywordReserved     InlineStyle `mapstructure:"keyword_reserved"`
	KeywordNamespace    InlineStyle `mapstructure:"keyword_namespace"`
	KeywordType         InlineStyle `mapstructure:"keyword_type"`
	Operator            InlineStyle `mapstructure:"operator"`
	Punctuation         InlineStyle `mapstructure:"punctuation"`
	Name                InlineStyle `mapstructure:"name"`
	NameBuiltin         InlineStyle `mapstructure:"name_builtin"`
	NameTag             InlineStyle `mapstructure:"name_tag"`
	NameAttribute       InlineStyle `mapstructure:"name_attribute"`
	NameClass           InlineStyle `mapstructure:"name_class"`
	NameConstant        InlineStyle `mapstructure:"name_constant"`
	NameDecorator       InlineStyle `mapstructure:"name_decorator"`
	NameException       InlineStyle `mapstructure:"name_exception"`
	NameFunction        InlineStyle `mapstructure:"name_function"`
	NameOther           InlineStyle `mapstructure:"name_other"`
	Literal             InlineStyle `mapstructure:"literal"`
	LiteralNumber       InlineStyle `mapstructure:"literal_number"`
	LiteralDate         InlineStyle `mapstructure:"literal_date"`
	LiteralString       InlineStyle `mapstructure:"literal_string"`
	LiteralStringEscape InlineStyle `mapstructure:"literal_string_escape"`
	GenericDeleted      InlineStyle `mapstructure:"generic_deleted"`
	GenericEmph         InlineStyle `mapstructure:"generic_emph"`
	GenericInserted     InlineStyle `mapstructure:"generic_inserted"`
	GenericStrong       InlineStyle `mapstructure:"generic_strong"`
	GenericSubheading   InlineStyle `mapstructure:"generic_subheading"`
	Background          InlineStyle `mapstructure:"background"`
}

type TableStyle struct {
	BlockStyle      `mapstructure:",squash"`
	Header          BlockStyle `mapstructure:"header"`
	Cell            BlockStyle `mapstructure:"cell"`
	CenterSeparator string     `mapstructure:"center_separator"`
	ColumnSeparator string     `mapstructure:"column_separator"`
	RowSeparator    string     `mapstructure:"row_separator"`
}

type TaskStylesConfig struct {
	FocusedBar     string         `mapstructure:"focused_color"`
	UnfocusedBar   string         `mapstructure:"unfocused_color"`
	CompletedStyle TaskStateStyle `mapstructure:"completed_style"`
	PendingStyle   TaskStateStyle `mapstructure:"pending_style"`
	AbandonedStyle TaskStateStyle `mapstructure:"abandoned_style"`
	StartedStyle   TaskStateStyle `mapstructure:"started_style"`
}

type TaskStateStyle struct {
	Title string `mapstructure:"color"`
	Desc  string `mapstructure:"text_color"`
}
