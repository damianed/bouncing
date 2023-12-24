package main

import (
	"fmt"
	"runtime"
  "log"
  "strings"

	"github.com/go-gl/gl/all-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
  triangle = []float32{
    -0.5, 0.5, 0, // top
    -0.5, -0.5, 0, // left
    0.5, -0.5, 0, // right
  }

  triangle2 = []float32{
    -0.5, 0.5, 0, // top
    0.5, 0.5, 0, // left
    0.5, -0.5, 0, // right
  }
)

const (
  width = 500
  height = 500
  vertexShaderSource = `
    #version 410
    in vec3 vp;
    void main() {
        gl_Position = vec4(vp, 1.0);
    }
  ` + "\x00"

  fragmentShaderSource = `
      #version 410
      layout(location = 0) out vec4 red_colour;
      layout(location = 1) out vec4 blue_colour;
      void main() {
          blue_colour = vec4(0, 0.2, 0.5, 1);
          red_colour = vec4(1, 0, 0, 1);
      }
  ` + "\x00"
)

func main() {
  runtime.LockOSThread()

  window := initGLFW()
  defer glfw.Terminate()

  program := initOpenGL()

  var vaoList []uint32
  vaoList = append(vaoList, makeVao(triangle, 0))
  vaoList = append(vaoList, makeVao(triangle2, 1))
  for !window.ShouldClose() {
    draw(vaoList, window, program)
  }
}

func initGLFW() (window *glfw.Window) {
  if err := glfw.Init(); err != nil {
    panic(err)
  }

  glfw.WindowHint(glfw.Resizable, glfw.False)
  glfw.WindowHint(glfw.ContextVersionMajor, 4)
  glfw.WindowHint(glfw.ContextVersionMinor, 1)
  glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
  glfw.WindowHint(glfw.OpenGLForwardCompatible, gl.TRUE)

  window, err := glfw.CreateWindow(width, height, "Golang OpenGL Triangle", nil, nil)
  if err != nil {
    panic(err)
  }
  window.MakeContextCurrent()

  return window
}

func initGL(window *glfw.Window) (uint32, error) {
  if err := gl.Init(); err != nil {
    panic(err)
  }

  version := gl.GoStr(gl.GetString(gl.VERSION))
  glsl := gl.GoStr(gl.GetString(gl.SHADING_LANGUAGE_VERSION))
  fmt.Println("OpenGL version", version, glsl)

  width, height := window.GetFramebufferSize()
  gl.Viewport(0, 0, int32(width), int32(height))

  prog := gl.CreateProgram()
  gl.LinkProgram(prog)

  return prog, nil
}

func initOpenGL() uint32 {
    if err := gl.Init(); err != nil {
            panic(err)
    }
    version := gl.GoStr(gl.GetString(gl.VERSION))
    log.Println("OpenGL version", version)

    vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
    if err != nil {
      panic(err)
    }

    fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
    if err != nil {
      panic(err)
    }

    //fragmentShader2, err := compileShader(fragmentShaderSource2, gl.FRAGMENT_SHADER)
    //if err != nil {
    //  panic(err)
    //}

    prog := gl.CreateProgram()
    gl.AttachShader(prog, vertexShader)
    gl.AttachShader(prog, fragmentShader)
    gl.LinkProgram(prog)
    return prog
}

func draw(vaoList []uint32, window *glfw.Window, prog uint32) {
  gl.Clear(gl.COLOR_BUFFER_BIT | gl.DEPTH_BUFFER_BIT)
  gl.UseProgram(prog)

  for index, vao := range vaoList {
    ui8 := uint8(index)
    gl.BindFragDataLocation(prog, uint32(index), &ui8)
    gl.BindVertexArray(vao)
    gl.DrawArrays(gl.TRIANGLES, 0, int32(len(triangle) / 3))
  }

  glfw.PollEvents()
  window.SwapBuffers()
}

func makeVao(points []float32, index uint32) uint32 {
  var vbo uint32
  gl.GenBuffers(1, &vbo)
  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.BufferData(gl.ARRAY_BUFFER, 4*len(points), gl.Ptr(points), gl.STATIC_DRAW)

  var vao uint32
  gl.GenVertexArrays(1, &vao)
  gl.BindVertexArray(vao)
  gl.EnableVertexAttribArray(index)
  gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
  gl.VertexAttribPointer(index, 3, gl.FLOAT, false, 0, nil)

  return vao
}


func compileShader(source string, shaderType uint32) (uint32, error) {
    shader := gl.CreateShader(shaderType)

    csources, free := gl.Strs(source)
    gl.ShaderSource(shader, 1, csources, nil)
    free()
    gl.CompileShader(shader)

    var status int32
    gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
    if status == gl.FALSE {
        var logLength int32
        gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)

        log := strings.Repeat("\x00", int(logLength+1))
        gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(log))

        return 0, fmt.Errorf("failed to compile %v: %v", source, log)
    }

    return shader, nil
}
