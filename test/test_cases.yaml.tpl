{{ range $key, $value := .}}
-
  id: "{{ $key }}"
  input: {{ $value.input }}
{{ end }}
