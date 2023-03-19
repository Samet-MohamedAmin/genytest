{{ range $key, $value := .}}
-
  id: "{{/* $key */}}"
  input:
    -
      name: {{ $value.name }}
      sellIn: {{ $value.sellIn }}
      quality: {{ $value.quality }}
{{ end }}
