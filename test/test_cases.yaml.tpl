{{ range $value := .}}
-
  comboHash: "{{ $value.Hash }}"
  input: {{ $value.Items.input }}
  bla: {{ $value.Items.bla }}
{{ end }}
