-- ============================================================
-- Novel2Script-AI 种子数据脚本
-- 用于本地开发环境的测试数据
-- ============================================================

USE novel2script;

-- ============================================================
-- 1. 测试用户
-- 密码均为 bcrypt hash of "test123456"
-- ============================================================
INSERT INTO sys_user (username, email, password_hash, role, status, created_at, updated_at) VALUES
('admin', 'admin@novel2script.local', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'admin', 'active', NOW(), NOW()),
('testuser1', 'test1@novel2script.local', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', 'active', NOW(), NOW()),
('testuser2', 'test2@novel2script.local', '$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy', 'user', 'active', NOW(), NOW())
ON DUPLICATE KEY UPDATE username=username;

-- ============================================================
-- 2. 测试项目
-- ============================================================
INSERT INTO novel_project (owner_id, title, description, adaptation_mode, visibility, status, created_at, updated_at) VALUES
(2, '雨夜重逢', '一个关于旧友重逢的短篇小说改编项目', 'screen_script', 'private', 'created', NOW(), NOW()),
(2, '江湖旧事', '武侠小说改编，讲述江湖恩怨', 'screen_script', 'private', 'uploaded', NOW(), NOW()),
(3, '都市传说', '都市悬疑短篇改编', 'short_video', 'private', 'created', NOW(), NOW())
ON DUPLICATE KEY UPDATE title=title;

-- ============================================================
-- 3. 测试章节
-- ============================================================
INSERT INTO novel_chapter (project_id, chapter_index, chapter_title, content, content_hash, created_at) VALUES
(1, 1, '第一章 重逢',
'林舟推开旧教学楼的门，走廊里弥漫着潮湿的气息。十年了，这里几乎没有变化。

他走到三楼的教室门前，透过玻璃窗看到一个熟悉的身影。苏晚坐在窗边，夕阳落在她的肩上，她没有回头。

"你终于来了。"她的声音很轻，却清晰地传入他的耳朵。

林舟推开门，脚步停在门口。他想说些什么，却发现喉咙像被什么堵住了一样。

"对不起，我来晚了。"他终于开口，声音有些沙哑。',
SHA2('chapter1_content', 256), NOW()),

(1, 2, '第二章 暗流',
'第二天，林舟回到学校，发现办公室的桌上放着一封没有署名的信。

信封里只有一张纸条：「当年的事，真相不是你想的那样。」

他的手微微颤抖。十年前的那个雨夜，他选择了离开，以为这样就能结束一切。

但现在，似乎有人不想让他忘记。',
SHA2('chapter2_content', 256), NOW()),

(1, 3, '第三章 雨夜',
'暴雨如注，林舟站在旧教学楼前。手机屏幕上是苏晚发来的消息：「我在老地方等你。」

他犹豫了很久，最终还是撑起伞走了出去。

当他到达教室时，苏晚已经在那里了。她的衣服被雨水打湿，但她似乎并不在意。

"我找到了当年的监控录像。"她说，声音平静得可怕。',
SHA2('chapter3_content', 256), NOW())
ON DUPLICATE KEY UPDATE chapter_title=chapter_title;

-- ============================================================
-- 4. 测试人物档案
-- ============================================================
INSERT INTO character_profile (project_id, character_key, name, aliases, role_type, description, personality, relationships, source_refs, confidence, created_at, updated_at) VALUES
(1, 'char_001', '林舟', '["阿舟"]', 'protagonist', '性格冷静的男主角，大学教授，十年前离开故乡', '["冷静", "克制", "理性"]', '[{"target": "char_002", "relation": "青梅竹马"}]', '[{"chapter_index": 1, "paragraph_start": 1, "paragraph_end": 3}]', 95.00, NOW(), NOW()),
(1, 'char_002', '苏晚', '["小晚"]', 'protagonist', '女主角，林舟的旧友，一直在等待真相', '["坚韧", "执着", "温柔"]', '[{"target": "char_001", "relation": "青梅竹马"}]', '[{"chapter_index": 1, "paragraph_start": 4, "paragraph_end": 6}]', 95.00, NOW(), NOW()),
(1, 'char_003', '周衡', '["老周"]', 'supporting', '林舟的大学同学，知晓当年真相', '["正直", "沉默"]', '[{"target": "char_001", "relation": "同学"}]', '[{"chapter_index": 2, "paragraph_start": 1, "paragraph_end": 3}]', 80.00, NOW(), NOW())
ON DUPLICATE KEY UPDATE name=name;

-- ============================================================
-- 5. 测试剧情事件
-- ============================================================
INSERT INTO plot_event (project_id, event_key, chapter_index, trigger_text, action_text, result_text, importance, source_refs, confidence, created_at) VALUES
(1, 'event_001', 1, '林舟收到匿名信', '前往旧教学楼', '与苏晚重逢', 'high', '[{"chapter_index": 1, "paragraph_start": 1, "paragraph_end": 10}]', 90.00, NOW()),
(1, 'event_002', 2, '林舟发现匿名信', '阅读纸条内容', '开始怀疑当年的真相', 'high', '[{"chapter_index": 2, "paragraph_start": 1, "paragraph_end": 5}]', 85.00, NOW()),
(1, 'event_003', 3, '苏晚发来消息', '林舟冒雨前往', '苏晚找到监控录像', 'high', '[{"chapter_index": 3, "paragraph_start": 1, "paragraph_end": 8}]', 90.00, NOW())
ON DUPLICATE KEY UPDATE event_key=event_key;

-- ============================================================
-- 6. 审计日志示例
-- ============================================================
INSERT INTO audit_log (user_id, project_id, action, resource_type, resource_id, ip_address, request_id, created_at) VALUES
(2, 1, 'project.create', 'novel_project', 1, '127.0.0.1', 'req_seed_001', NOW()),
(2, 1, 'file.upload', 'novel_source_file', 1, '127.0.0.1', 'req_seed_002', NOW()),
(2, 1, 'task.create', 'ai_task', 1, '127.0.0.1', 'req_seed_003', NOW())
ON DUPLICATE KEY UPDATE action=action;
