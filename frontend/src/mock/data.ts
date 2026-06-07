// Mock 数据 —— Novel2Script-AI

export interface Project {
  id: number
  title: string
  description: string
  adaptation_mode: 'screen_script' | 'stage_play' | 'short_video' | 'radio_drama'
  status: 'created' | 'uploaded' | 'processing' | 'completed' | 'has_risk'
  chapter_count: number
  updated_at: string
  created_at: string
}

export interface Chapter {
  id: number
  project_id: number
  chapter_index: number
  chapter_title: string
  content: string
  word_count: number
  paragraph_count: number
}

export interface Character {
  id: string
  name: string
  aliases: string[]
  role: 'protagonist' | 'antagonist' | 'supporting' | 'minor'
  description: string
  personality: string[]
  relationships: { target: string; relation: string }[]
  source_trace: { chapter_index: number; paragraph_start: number; paragraph_end: number }[]
  confidence: number
}

export interface PlotEvent {
  id: string
  chapter_index: number
  trigger: string
  action: string
  result: string
  importance: 'high' | 'medium' | 'low'
  characters_involved: string[]
  source_trace: { chapter_index: number; paragraph_start: number; paragraph_end: number }
  confidence: number
}

export interface Scene {
  id: string
  title: string
  order: number
  time: string
  location: string
  characters: string[]
  summary: string
  source_trace: { chapter_index: number; paragraph_start: number; paragraph_end: number }[]
  risk_level: 'none' | 'low' | 'medium' | 'high'
  risk_type?: string
}

export interface ValidationIssue {
  id: number
  issue_type: string
  severity: 'low' | 'medium' | 'high'
  message: string
  location_path: string
  suggestion: string
  resolved: boolean
}

export interface Version {
  id: number
  version_no: number
  created_by: 'ai' | 'user' | 'system_repair'
  created_at: string
  validation_status: 'valid' | 'invalid' | 'pending'
  hallucination_risk: 'low' | 'medium' | 'high'
  safety_risk: 'low' | 'medium' | 'high'
}

export interface AuditLog {
  id: number
  action: string
  project_title: string
  user: string
  status: string
  request_id: string
  created_at: string
}

// ============ 项目 ============
export const mockProjects: Project[] = [
  {
    id: 1,
    title: '旧城往事',
    description: '一部关于江南小镇三代人命运纠葛的长篇小说',
    adaptation_mode: 'screen_script',
    status: 'completed',
    chapter_count: 3,
    updated_at: '2026-06-06 14:30',
    created_at: '2026-06-05 10:00'
  },
  {
    id: 2,
    title: '星际漫游指南',
    description: '科幻短篇集，讲述人类在银河系边缘的探索故事',
    adaptation_mode: 'short_video',
    status: 'uploaded',
    chapter_count: 0,
    updated_at: '2026-06-06 09:15',
    created_at: '2026-06-06 09:00'
  }
]

// ============ 章节 ============
export const mockChapters: Chapter[] = [
  {
    id: 1,
    project_id: 1,
    chapter_index: 1,
    chapter_title: '第一章 重逢',
    content: '林舟推开旧教学楼的门，走廊里弥漫着潮湿的气息。十年了，这里几乎没有变化。\n\n他走到三楼的教室门前，透过玻璃窗看到一个熟悉的身影。苏晚坐在窗边，夕阳落在她的肩上，她没有回头。\n\n"你终于来了。"她的声音很轻，却清晰地传入他的耳朵。\n\n林舟推开门，脚步停在门口。他想说些什么，却发现喉咙像被什么堵住了一样。\n\n"对不起，我来晚了。"他终于开口，声音有些沙哑。',
    word_count: 156,
    paragraph_count: 7
  },
  {
    id: 2,
    project_id: 1,
    chapter_index: 2,
    chapter_title: '第二章 暗流',
    content: '第二天，林舟回到学校，发现办公室的桌上放着一封没有署名的信。\n\n信封里只有一张纸条：「当年的事，真相不是你想的那样。」\n\n他的手微微颤抖。十年前的那个雨夜，他选择了离开，以为这样就能结束一切。\n\n但现在，似乎有人不想让他忘记。',
    word_count: 98,
    paragraph_count: 5
  },
  {
    id: 3,
    project_id: 1,
    chapter_index: 3,
    chapter_title: '第三章 雨夜',
    content: '暴雨如注，林舟站在旧教学楼前。手机屏幕上是苏晚发来的消息：「我在老地方等你。」\n\n他犹豫了很久，最终还是撑起伞走了出去。\n\n当他到达教室时，苏晚已经在那里了。她的衣服被雨水打湿，但她似乎并不在意。\n\n"我找到了当年的监控录像。"她说，声音平静得可怕。',
    word_count: 112,
    paragraph_count: 5
  }
]

// ============ 人物 ============
export const mockCharacters: Character[] = [
  {
    id: 'char_001',
    name: '林舟',
    aliases: ['阿舟'],
    role: 'protagonist',
    description: '十年前离开小镇的青年，如今回到旧地，试图寻找当年的真相',
    personality: ['内敛', '执着', '愧疚'],
    relationships: [
      { target: 'char_002', relation: '旧友，曾有未说出口的感情' },
      { target: 'char_003', relation: '恩师' }
    ],
    source_trace: [{ chapter_index: 1, paragraph_start: 1, paragraph_end: 3 }],
    confidence: 0.95
  },
  {
    id: 'char_002',
    name: '苏晚',
    aliases: ['晚晚'],
    role: 'protagonist',
    description: '留在小镇的女教师，一直守护着某个秘密',
    personality: ['冷静', '坚定', '隐忍'],
    relationships: [
      { target: 'char_001', relation: '旧友，一直在等他回来' },
      { target: 'char_004', relation: '同事' }
    ],
    source_trace: [{ chapter_index: 1, paragraph_start: 3, paragraph_end: 4 }],
    confidence: 0.95
  },
  {
    id: 'char_003',
    name: '陈老师',
    aliases: ['陈老'],
    role: 'supporting',
    description: '退休老教师，知道当年事件的关键线索',
    personality: ['慈祥', '守口如瓶'],
    relationships: [
      { target: 'char_001', relation: '学生' },
      { target: 'char_002', relation: '学生' }
    ],
    source_trace: [{ chapter_index: 2, paragraph_start: 1, paragraph_end: 2 }],
    confidence: 0.80
  },
  {
    id: 'char_004',
    name: '方远',
    aliases: [],
    role: 'supporting',
    description: '学校新来的体育老师，对苏晚有好感',
    personality: ['直率', '热心'],
    relationships: [
      { target: 'char_002', relation: '同事，暗恋' }
    ],
    source_trace: [{ chapter_index: 2, paragraph_start: 3, paragraph_end: 4 }],
    confidence: 0.70
  },
  {
    id: 'char_005',
    name: '匿名信作者',
    aliases: ['神秘人'],
    role: 'minor',
    description: '身份不明，寄出匿名信的人',
    personality: ['未知'],
    relationships: [],
    source_trace: [{ chapter_index: 2, paragraph_start: 1, paragraph_end: 2 }],
    confidence: 0.50
  }
]

// ============ 剧情事件 ============
export const mockPlotEvents: PlotEvent[] = [
  {
    id: 'event_001',
    chapter_index: 1,
    trigger: '林舟收到消息，决定回到小镇',
    action: '推开旧教学楼的门，走上三楼',
    result: '在教室里与苏晚重逢',
    importance: 'high',
    characters_involved: ['char_001', 'char_002'],
    source_trace: { chapter_index: 1, paragraph_start: 1, paragraph_end: 4 },
    confidence: 0.95
  },
  {
    id: 'event_002',
    chapter_index: 1,
    trigger: '苏晚说出"你终于来了"',
    action: '林舟道歉',
    result: '两人之间的紧张气氛有所缓和',
    importance: 'medium',
    characters_involved: ['char_001', 'char_002'],
    source_trace: { chapter_index: 1, paragraph_start: 4, paragraph_end: 7 },
    confidence: 0.90
  },
  {
    id: 'event_003',
    chapter_index: 2,
    trigger: '林舟回到学校办公室',
    action: '发现一封没有署名的信',
    result: '信中暗示当年事件另有隐情',
    importance: 'high',
    characters_involved: ['char_001'],
    source_trace: { chapter_index: 2, paragraph_start: 1, paragraph_end: 2 },
    confidence: 0.95
  },
  {
    id: 'event_004',
    chapter_index: 2,
    trigger: '林舟读到纸条内容',
    action: '回忆十年前的雨夜',
    result: '决定追查真相',
    importance: 'high',
    characters_involved: ['char_001'],
    source_trace: { chapter_index: 2, paragraph_start: 2, paragraph_end: 4 },
    confidence: 0.85
  },
  {
    id: 'event_005',
    chapter_index: 3,
    trigger: '暴雨夜收到苏晚消息',
    action: '林舟撑伞前往旧教学楼',
    result: '与苏晚在教室见面',
    importance: 'medium',
    characters_involved: ['char_001', 'char_002'],
    source_trace: { chapter_index: 3, paragraph_start: 1, paragraph_end: 3 },
    confidence: 0.90
  },
  {
    id: 'event_006',
    chapter_index: 3,
    trigger: '苏晚提到监控录像',
    action: '展示找到的证据',
    result: '真相开始浮出水面',
    importance: 'high',
    characters_involved: ['char_001', 'char_002'],
    source_trace: { chapter_index: 3, paragraph_start: 4, paragraph_end: 5 },
    confidence: 0.85
  },
  {
    id: 'event_007',
    chapter_index: 1,
    trigger: '林舟站在教室门口',
    action: '犹豫是否进入',
    result: '内心的愧疚与期待交织',
    importance: 'low',
    characters_involved: ['char_001'],
    source_trace: { chapter_index: 1, paragraph_start: 5, paragraph_end: 6 },
    confidence: 0.80
  },
  {
    id: 'event_008',
    chapter_index: 2,
    trigger: '林舟回忆离开的决定',
    action: '反思过去的选择',
    result: '意识到逃避并不能解决问题',
    importance: 'medium',
    characters_involved: ['char_001'],
    source_trace: { chapter_index: 2, paragraph_start: 3, paragraph_end: 5 },
    confidence: 0.75
  }
]

// ============ 场景 ============
export const mockScenes: Scene[] = [
  {
    id: 'scene_001',
    title: '旧教室重逢',
    order: 1,
    time: '傍晚',
    location: '旧教学楼三楼教室',
    characters: ['char_001', 'char_002'],
    summary: '林舟推开旧教学楼的门，在三楼教室与苏晚重逢。苏晚坐在窗边，夕阳落在她的肩上。',
    source_trace: [{ chapter_index: 1, paragraph_start: 1, paragraph_end: 5 }],
    risk_level: 'none'
  },
  {
    id: 'scene_002',
    title: '久别重逢的对话',
    order: 2,
    time: '傍晚',
    location: '旧教学楼三楼教室',
    characters: ['char_001', 'char_002'],
    summary: '林舟向苏晚道歉，两人之间的紧张气氛有所缓和。',
    source_trace: [{ chapter_index: 1, paragraph_start: 5, paragraph_end: 7 }],
    risk_level: 'none'
  },
  {
    id: 'scene_003',
    title: '匿名信',
    order: 3,
    time: '白天',
    location: '学校办公室',
    characters: ['char_001'],
    summary: '林舟在办公室发现一封匿名信，暗示当年事件另有隐情。',
    source_trace: [{ chapter_index: 2, paragraph_start: 1, paragraph_end: 3 }],
    risk_level: 'none'
  },
  {
    id: 'scene_004',
    title: '十年前的回忆',
    order: 4,
    time: '白天',
    location: '学校办公室',
    characters: ['char_001'],
    summary: '林舟读到纸条，回忆十年前的雨夜，决定追查真相。',
    source_trace: [{ chapter_index: 2, paragraph_start: 3, paragraph_end: 5 }],
    risk_level: 'low',
    risk_type: '心理描写需转化为动作'
  },
  {
    id: 'scene_005',
    title: '暴雨夜赴约',
    order: 5,
    time: '深夜',
    location: '旧教学楼',
    characters: ['char_001'],
    summary: '暴雨夜，林舟收到苏晚消息后撑伞前往旧教学楼。',
    source_trace: [{ chapter_index: 3, paragraph_start: 1, paragraph_end: 3 }],
    risk_level: 'none'
  },
  {
    id: 'scene_006',
    title: '监控录像',
    order: 6,
    time: '深夜',
    location: '旧教学楼三楼教室',
    characters: ['char_001', 'char_002'],
    summary: '苏晚展示找到的监控录像，真相开始浮出水面。',
    source_trace: [{ chapter_index: 3, paragraph_start: 4, paragraph_end: 5 }],
    risk_level: 'medium',
    risk_type: '关键情节缺少细节展开'
  }
]

// ============ YAML 剧本 ============
export const mockYamlScript = `script:
  metadata:
    title: 旧城往事
    source_title: 旧城往事
    adaptation_mode: screen_script
    language: zh-CN
    version: 1

  characters:
    - id: char_001
      name: 林舟
      aliases:
        - 阿舟
      role: protagonist
      description: 十年前离开小镇的青年，如今回到旧地，试图寻找当年的真相
      personality:
        - 内敛
        - 执着
        - 愧疚
      relationships:
        - target: char_002
          relation: 旧友，曾有未说出口的感情
        - target: char_003
          relation: 恩师
      source_trace:
        - chapter_index: 1
          paragraph_start: 1
          paragraph_end: 3

    - id: char_002
      name: 苏晚
      aliases:
        - 晚晚
      role: protagonist
      description: 留在小镇的女教师，一直守护着某个秘密
      personality:
        - 冷静
        - 坚定
        - 隐忍
      relationships:
        - target: char_001
          relation: 旧友，一直在等他回来
      source_trace:
        - chapter_index: 1
          paragraph_start: 3
          paragraph_end: 4

  scenes:
    - id: scene_001
      title: 旧教室重逢
      order: 1
      time: 傍晚
      location: 旧教学楼三楼教室
      characters:
        - char_001
        - char_002
      summary: 林舟推开旧教学楼的门，在三楼教室与苏晚重逢。
      source_trace:
        - chapter_index: 1
          paragraph_start: 1
          paragraph_end: 5
      actions:
        - character: char_001
          description: 推开教室门，脚步停在门口
        - character: char_002
          description: 坐在窗边，没有回头
      dialogues:
        - character: char_002
          line: 你终于来了。
          emotion: 平静
        - character: char_001
          line: 对不起，我来晚了。
          emotion: 愧疚

    - id: scene_002
      title: 久别重逢的对话
      order: 2
      time: 傍晚
      location: 旧教学楼三楼教室
      characters:
        - char_001
        - char_002
      summary: 林舟向苏晚道歉，两人之间的紧张气氛有所缓和。
      source_trace:
        - chapter_index: 1
          paragraph_start: 5
          paragraph_end: 7
      dialogues:
        - character: char_001
          line: 这些年，你过得好吗？
          emotion: 试探

    - id: scene_003
      title: 匿名信
      order: 3
      time: 白天
      location: 学校办公室
      characters:
        - char_001
      summary: 林舟在办公室发现一封匿名信，暗示当年事件另有隐情。
      source_trace:
        - chapter_index: 2
          paragraph_start: 1
          paragraph_end: 3
      actions:
        - character: char_001
          description: 发现桌上的匿名信，拆开阅读
      adaptation_notes:
        - source: 原文第2章第1段
          change: 增加拆信动作描写
          reason: 剧本需要可视化动作
`

// ============ 校验问题 ============
export const mockValidationIssues: ValidationIssue[] = [
  {
    id: 1,
    issue_type: 'missing_source_trace',
    severity: 'medium',
    message: '场景 scene_002 缺少完整的原文溯源信息',
    location_path: 'script.scenes[1].source_trace',
    suggestion: '补充第二章第5-7段的段落引用',
    resolved: false
  },
  {
    id: 2,
    issue_type: 'character_reference',
    severity: 'low',
    message: '人物 char_003（陈老师）在剧本中被引用但未在人物档案中定义',
    location_path: 'script.characters',
    suggestion: '在人物档案中添加陈老师的完整信息',
    resolved: false
  },
  {
    id: 3,
    issue_type: 'dialogue_style',
    severity: 'low',
    message: '场景 scene_002 中林舟的对白风格与人物设定略有偏差',
    location_path: 'script.scenes[1].dialogues[0]',
    suggestion: '林舟性格内敛，对白可更含蓄',
    resolved: true
  },
  {
    id: 4,
    issue_type: 'scene_summary',
    severity: 'low',
    message: '场景 scene_004 的摘要过于简短',
    location_path: 'script.scenes[3].summary',
    suggestion: '补充更多场景细节描述',
    resolved: false
  },
  {
    id: 5,
    issue_type: 'format_error',
    severity: 'high',
    message: 'YAML 缩进错误：scene_006 的 actions 字段缩进不一致',
    location_path: 'script.scenes[5].actions',
    suggestion: '统一使用 2 空格缩进',
    resolved: false
  }
]

// ============ 版本记录 ============
export const mockVersions: Version[] = [
  {
    id: 1,
    version_no: 3,
    created_by: 'user',
    created_at: '2026-06-06 14:30',
    validation_status: 'valid',
    hallucination_risk: 'low',
    safety_risk: 'low'
  },
  {
    id: 2,
    version_no: 2,
    created_by: 'system_repair',
    created_at: '2026-06-06 13:15',
    validation_status: 'valid',
    hallucination_risk: 'low',
    safety_risk: 'low'
  },
  {
    id: 3,
    version_no: 1,
    created_by: 'ai',
    created_at: '2026-06-06 12:00',
    validation_status: 'invalid',
    hallucination_risk: 'medium',
    safety_risk: 'low'
  },
  {
    id: 4,
    version_no: 0,
    created_by: 'ai',
    created_at: '2026-06-06 11:30',
    validation_status: 'invalid',
    hallucination_risk: 'high',
    safety_risk: 'low'
  }
]

// ============ 审计日志 ============
export const mockAuditLogs: AuditLog[] = [
  { id: 1, action: 'user.login', project_title: '-', user: 'testuser', status: '成功', request_id: 'req_a1b2c3', created_at: '2026-06-06 14:30:00' },
  { id: 2, action: 'project.create', project_title: '旧城往事', user: 'testuser', status: '成功', request_id: 'req_d4e5f6', created_at: '2026-06-06 14:31:00' },
  { id: 3, action: 'file.upload', project_title: '旧城往事', user: 'testuser', status: '成功', request_id: 'req_g7h8i9', created_at: '2026-06-06 14:32:00' },
  { id: 4, action: 'task.create', project_title: '旧城往事', user: 'testuser', status: '成功', request_id: 'req_j0k1l2', created_at: '2026-06-06 14:33:00' },
  { id: 5, action: 'script.generate', project_title: '旧城往事', user: 'system', status: '成功', request_id: 'req_m3n4o5', created_at: '2026-06-06 14:35:00' },
  { id: 6, action: 'script.validate', project_title: '旧城往事', user: 'testuser', status: '有警告', request_id: 'req_p6q7r8', created_at: '2026-06-06 14:36:00' },
  { id: 7, action: 'script.edit', project_title: '旧城往事', user: 'testuser', status: '成功', request_id: 'req_s9t0u1', created_at: '2026-06-06 14:37:00' },
  { id: 8, action: 'script.export', project_title: '旧城往事', user: 'testuser', status: '成功', request_id: 'req_v2w3x4', created_at: '2026-06-06 14:38:00' }
]

// ============ 工作流步骤 ============
export interface WorkflowStep {
  key: string
  label: string
  status: 'pending' | 'running' | 'completed' | 'warning' | 'failed'
  count?: number
  unit?: string
  detail?: string
}

export const mockWorkflowSteps: WorkflowStep[] = [
  { key: 'clean', label: '文本清洗', status: 'completed', detail: '去除 12 处乱码，规范化空白字符' },
  { key: 'split', label: '章节切分', status: 'completed', count: 3, unit: '章', detail: '识别到 3 个章节标题' },
  { key: 'character', label: '人物抽取', status: 'completed', count: 5, unit: '人', detail: '提取 5 个人物档案' },
  { key: 'plot', label: '剧情事件链', status: 'completed', count: 8, unit: '个', detail: '构建 8 个剧情事件' },
  { key: 'scene', label: '场景拆分', status: 'completed', count: 6, unit: '场', detail: '拆分为 6 个场景' },
  { key: 'generate', label: '剧本生成', status: 'completed', detail: '生成 YAML 剧本，3 个场景含对白' },
  { key: 'validate', label: 'Schema 校验', status: 'warning', detail: '发现 5 个校验问题' },
  { key: 'hallucination', label: '幻觉检测', status: 'completed', detail: '风险等级：低' },
  { key: 'safety', label: '安全审查', status: 'completed', detail: '未发现安全风险' }
]

// ============ 改编模式说明 ============
export const adaptationModes = {
  screen_script: { label: '影视剧本', desc: '适用于电影、电视剧改编，注重场景转换和镜头语言' },
  stage_play: { label: '舞台剧', desc: '适用于话剧改编，注重对白张力和舞台调度' },
  short_video: { label: '短视频分镜', desc: '适用于短剧改编，注重节奏和视觉冲击' },
  radio_drama: { label: '广播剧', desc: '适用于音频改编，注重声音表现和听觉节奏' }
}
