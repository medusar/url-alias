package main

import (
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"html/template"
	"net/http"
)

type PageValue struct {
	List []*Alias
}

type Alias struct {
	Alias string
	Url   string
}

func main() {
	var templates = template.Must(template.ParseFiles("html/config.html"))

	initConfig()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path[1:]
		if url := viper.GetString(path); url != "" {
			http.Redirect(w, r, url, http.StatusFound)
		} else {
			http.Redirect(w, r, "/config", http.StatusFound)
		}
	})

	// config html
	http.HandleFunc("/config", func(writer http.ResponseWriter, r *http.Request) {
		list := make([]*Alias, 0)
		for _, key := range viper.AllKeys() {
			if value := viper.GetString(key); value != "" {
				list = append(list, &Alias{key, value})
			}
		}
		templates.ExecuteTemplate(writer, "config.html", &PageValue{List: list})
	})

	http.ListenAndServe(":80", nil)
}

func initConfig() {
	viper.SetConfigName("config")           // name of config file (without extension)
	viper.SetConfigType("yaml")             // REQUIRED if the config file does not have the extension in the name
	viper.AddConfigPath("$HOME/.url-alias") // call multiple times to add many search paths
	viper.AddConfigPath(".")                // optionally look for config in the working directory
	err := viper.ReadInConfig()             // Find and read the config file
	if err != nil {                         // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		fmt.Println("Config file changed:", e.Name)
	})
}
