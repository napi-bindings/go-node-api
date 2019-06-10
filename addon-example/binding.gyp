{
  "targets": [
    {
      "target_name": "hello",
      "sources": [ "hello.cc" ],
      "include_dirs": [
        "include"
      ],
      "libraries": [
        "<(module_root_dir)/libgoaddon.a"
      ]
    }
  ]
}