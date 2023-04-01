<template>
  <el-container v-loading.fullscreen.lock="fullscreenLoading" style="width: 100%; height: 100%;">
    <el-container style="width: 100%; height: 100%;">
      <el-main>
        <div class="jsmind_layout">
            <el-icon @click="main"><Back /></el-icon>
          <div
            id="jsmind_container"
            ref="container"
            >
          </div>
          

        </div>

        <!-- TODO 拖动左右两侧div调整宽度 -->
        <!-- <div class="resize" title="收缩侧边栏">
          ⋮
        </div> -->

        <el-container class="answer_layout">
          <div v-html="answer">

          </div>
        </el-container>
      </el-main>
    
      <el-footer style="height: 25%;">
        <div style="width: 100%; height: 95%;">
          <el-card class="footer-card">
            <div>
            <el-form-item class="footer-card-button">
              
              <el-button type="primary" size="large" @click="addBrother">插入平级</el-button>
              <el-button type="primary" size="large" @click="addChild">插入子级</el-button>
              <el-button type="primary" size="large" @click="delCard">删除卡片</el-button>
              <el-button type="primary" size="large" @click="request">确定</el-button>
            </el-form-item>
            <el-form-item class="form">
                <el-select v-model="type" :disabled="disabled" @change="typechange" class="m-2" placeholder="请选择查询类别">
                    <el-option
                    v-for="item in types"
                    :key="item"
                    :label="item"
                    :value="item"
                    />
                </el-select>
                <el-select v-model="nextrule" :disabled="disablednextrule" class="m-2" placeholder="请选择类型">
                    <el-option
                    v-for="item in next_rule"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                    />
                </el-select>
                <el-select v-show="showlanguage" v-model="selectlanguage" class="m-2" placeholder="请选择程序语言">
                    <el-option
                    v-for="item in language"
                    :key="item.value"
                    :label="item.label"
                    :value="item.value"
                    />
                </el-select>
                
              </el-form-item>
              <el-form-item class="footer-card-texterea">
                <el-input v-model="content" :disabled="disabled" type="textarea" @keyup.enter="request" placeholder="在此输入问题，我是富文本,可以多行输入，balabalabala....."/>
              </el-form-item>
            </div>
          </el-card>
        </div>
      </el-footer>
    </el-container>
  </el-container>
</template>

<script>
import { ref } from 'vue'
import router from '@/router/index'
import 'jsmind/style/jsmind.css'
import jsMind from 'jsmind/js/jsmind.js'
import { marked } from 'marked'
import axios from "axios"

export default {
  watch: {
    
  },
data () {
return {
  fullscreenLoading: false,  // 全屏加载
  disabled: false,  // 禁止输入
  disablednextrule: false,
  showlanguage: false,
  answer: "",
  content: "",
  topic: "New Node",
  selectlanguage: "",
  language: [
    {
      value: 'golang',
      label: 'golang',
    },
    {
      value: 'python',
      label: 'python',
    },
    {
      value: 'powershell',
      label: 'powershell',
    },
    {
      value: 'c++',
      label: 'c++',
    },
    {
      value: 'cmd',
      label: 'cmd',
    },
  ],
  type: "",
  types: ([
    "chat",
    // "Code",
    // "ReadDoc"
  ]),
  nextrule: "",
  next_rule:[
    {
      value: "",
      label: "",
    },
  ],

  mind: {
    // 元数据
    meta: {
      name: '思维导图',
      author: 'effibot',
      version: '0.2'
    },
    // 数据格式声明
    format: 'node_tree',
    // 数据内容
    data: {
      id: 'root',
      topic: 'Root',
      answerid: '',
      answer: '',
    }
  },
  options: {
    container: 'jsmind_container', // [必选] 容器的ID
    editable: true, // [可选] 是否启用编辑
    theme: '', // [可选] 主题
    view: {
      engine: 'canvas', // 思维导图各节点之间线条的绘制引擎
      hmargin: 120, // 思维导图距容器外框的最小水平距离
      vmargin: 50, // 思维导图距容器外框的最小垂直距离
      line_width: 2, // 思维导图线条的粗细
      line_color: '#ddd' // 思维导图线条的颜色
    },
    layout: {
      hspace: 100, // 节点之间的水平间距
      vspace: 20, // 节点之间的垂直间距
      pspace: 20 // 节点与连接线之间的水平间距（用于容纳节点收缩/展开控制器）
    },
    shortcut: {
      enable: true, // 是否启用快捷键 默认为true
      handles: {
      },
      mapping: {
          addchild: [45, 4096+13], // Insert, Ctrl+Enter
          addbrother: 13, // Enter
          editnode: 113,// F2
          delnode: 46, // Delete
          toggle: 32, // Space
          left: 37, // Left
          up: 38, // Up
          right: 39, // Right
          down: 40, // Down
      }
    },
    zoom: {
      value: 100, // 层级大小
      min: 10, // 最小层级
      max: 400 // 最大层级
    },
  },
}
  
},

mounted () {
// 初始化
    this.jm = jsMind.show(this.options, this.mind)
    this.editor = this.jm.view.e_editor
    this.init()
    this.mouseWheel()
    this.mouseDrag()
},

methods: {
  // 初始化配置
  init () {
      this.jm.view.minZoom = 0.1
      this.jm.view.maxZoom = 5
      this.jm.expand_all()
      this.jm.select_node('root')
    },

  // 随机数生成节点ID
  nodeID () {
    return Math.random().toString(16).substring(2)
  },
  
  // 插入平级
  addBrother () {
    var nodeid = this.nodeID()
    const selectedNode = this.jm.get_selected_node()
    if (!!selectedNode && !selectedNode.isroot) {
      var node = this.jm.insert_node_after(selectedNode, nodeid, this.topic);
      if (!!node) {
        this.jm.select_node(nodeid);
      }
    }
    else {
      this.$message.error('请选择卡片')
    }
  },

  // 插入子级
  addChild () {
    var nodeid = this.nodeID()
    const selectedNode = this.jm.get_selected_node()
    if (!!selectedNode) {
      var node = this.jm.add_node(selectedNode, nodeid, this.topic);
      if (!!node) {
        this.jm.select_node(nodeid);
      }
    } else {
      this.$message.error('请选择卡片')
    }
  },

  // 删除卡片
  delCard () {
    var selectedNode = this.jm.get_selected_node()
    if (!!selectedNode && !selectedNode.isroot) {
      // TODO BUG:先添加卡片，然后改名，然后删除会报错
      this.jm.select_node(selectedNode.parent)
      this.jm.remove_node(selectedNode)
      // 获取数据
      console.log(this.jm.get_data('node_tree'))
    } else {
      this.$message.error('请选择卡片')
    }
  },

  // 缩小
  zoomOut () {
    this.jm.view.zoomOut()
  },

  // 放大
  zoomIn () {
    this.jm.view.zoomIn()
  },

  // 鼠标滚轮放大缩小
  mouseWheel () {
    if (document.addEventListener) {
      document.addEventListener('domMouseScroll', this.scrollFunc, false)
    }
    this.$refs.container.onmousewheel = this.scrollFunc
  },

  // 滚轮缩放
  scrollFunc (e) {
    e = e || window.event
    if (e.wheelDelta) {
      if (e.wheelDelta > 0) {
        this.zoomIn()
      } else {
        this.zoomOut()
      }
    } else if (e.detail) {
      if (e.detail > 0) {
        this.zoomIn()
      } else {
        this.zoomOut()
      }
    }
    e.preventDefault()
    this.jm.resize()
  },

  // 鼠标拖拽
  mouseDrag () {
    // 里层
    const el = document.querySelector('.jsmind-inner')
    // 选中节点
    let selected
    el.onmousedown = (ev) => {
      // 选中节点
      selected = this.jm.get_selected_node()
      // 标识 是否拖拽节点 避免冲突
      this.dragNodeFlag = !!selected
      const disX = ev.clientX
      const disY = ev.clientY
      const originalScrollLeft = el.scrollLeft
      const originalScrollTop = el.scrollTop
      const originalScrollBehavior = el.style['scroll-behavior']
      const originalPointerEvents = el.style['pointer-events']
      // auto: 默认值，表示滚动框立即滚动到指定位置。
      el.style['scroll-behavior'] = 'auto'
      // 鼠标移动事件是监听的整个document，这样可以使鼠标能够在元素外部移动的时候也能实现拖动
      document.onmousemove = (ev) => {
        if (this.dragNodeFlag) return
        this.drag = false
        ev.preventDefault()
        // 计算拖拽的偏移距离
        const distanceX = ev.clientX - disX
        const distanceY = ev.clientY - disY
        el.scrollTo(originalScrollLeft - distanceX, originalScrollTop - distanceY)
        // 在鼠标拖动的时候将点击事件屏蔽掉
        el.style['pointer-events'] = 'none'
        el.style.cursor = 'grabbing'
      }
      document.onmouseup = () => {
        if (!this.dragNodeFlag) {
          el.style['scroll-behavior'] = originalScrollBehavior
          el.style['pointer-events'] = originalPointerEvents
          el.style.cursor = 'grab'
        }
        document.onmousemove = document.onmouseup = null
      }
    }
  },

  // 上边是jsmind语法，下边是请求数据相关的方法

  main () {
    router.push("/")
  },

  typechange () {
    switch ( this.type ) {

      case "chat" :
        this.disablednextrule = true
        this.nextrule = "chat";
        this.next_rule = [
          {
            value: "chat",
            label: "chat",
          },
        ];
        return;
    }
  },

  // 请求后端数据
  request () {
    var selectedNode = this.jm.get_selected_node()

    var pid = ""
    if ( selectedNode === null ) {
      this.$message.error('请选择卡片');
      return;
    }
    if ( this.type === "" ) {
      this.$message.error('请选择查询类别');
      return;
    }
    if ( this.nextrule === "" ) {
      this.$message.error('请选择类型');
      return;
    }
    if ( selectedNode.answerid != null ) {
      this.answer = selectedNode.answer
      if ( this.content != "" ) {
        this.$message.error('此节点已经记录问题，请选择其他节点，返回当前节点记录问题')
        this.content = ""
      }
      return
    }
    if ( selectedNode.isroot ) {
      pid = ""
    } else {
      pid = selectedNode.parent.answerid
    }
    if ( !selectedNode.isroot ) {
      if ( selectedNode.parent != null && selectedNode.parent.answerid == null ) {
        this.$message.error('上一节点未记录问题');
        return
      }
    }
    if ( this.content === "" ) {
      this.$message.error('请输入你的问题');
      return;
    }
    this.jm.update_node(selectedNode.id, this.content)
    
    let url = '/api/answer?type_id=' + this.nextrule + '&content=' + this.content.replace(/;/g, '.').replace(/&/g, ' ').replace(/\?/g, '.') + '&p_id=' + pid + '&params=' + this.selectlanguage;
    this.fullscreenLoading = true
    axios.get(url)
      .then((data) => {
        this.fullscreenLoading = false
        if ( data.data.code === 200 ) {
          selectedNode.answerid = data.data.answer_id
          if ( selectedNode.isroot ) {
            this.answer = marked("\n## 问题：\n" + this.content + "<br /> <br /> \n## 回答：\n" + data.data.answer + "<br /> <br />")
          } else {
            this.answer = marked(selectedNode.parent.answer + "\n## 问题：\n" + this.content + "<br /> <br /> \n## 回答：\n" + data.data.answer + "<br /> <br />")
          }
          selectedNode.answer = this.answer
        } else {
          this.$message.error("非200返回码")
          console.log(data.data)
          return
        }
        this.content = ""
    })
      .catch((err) => {
        this.fullscreenLoading = false
        this.content = ""
        this.$message.error('err')
    })
  },
  

  }
}
</script>

<style lang="less" scoped>
.jsmind_layout {
  float: left;
  width: 70%;
  height: 100%;
}
.answer_layout {
  overflow: auto;
  width: 30%;
  height: 100%;
  float: right;
  box-shadow: -1px 4px 5px 3px rgba(0, 0, 0, 0.11);
  background-color: #F5F5F5;
}

#jsmind_container {
  width: 100%;
  height: 100%;
}
// .resize {
//   cursor: col-resize;
//   float: left;
//   position: relative;
//   top: 45%;
//   background-color: #d6d6d6;
//   border-radius: 5px;
//   margin-top: -10px;
//   width: 10px;
//   height: 50px;
//   background-size: cover;
//   background-position: center;
//   /*z-index: 99999;*/
//   font-size: 32px;
//   color: white;
//  }
//  // 拖拽区鼠标悬停样式
//  .resize:hover {
//   color: #444444;
//  }
.footer-card {
width: 100%;
height: 100%;
}
.footer-card-texterea {
  margin-bottom: 10px;
}

.footer-card-button {
  margin-left: 40%;
}
.form {
  margin-left: 40%;
}
</style>