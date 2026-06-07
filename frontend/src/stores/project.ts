import { defineStore } from 'pinia'
import { projectApi, chapterApi, scriptApi, auditApi, taskApi, resultApi } from '@/api'
import * as yaml from 'js-yaml'
import {
  mockCharacters,
  mockPlotEvents,
  mockScenes,
  mockYamlScript,
  mockValidationIssues,
  mockVersions,
  mockWorkflowSteps,
  type Project,
  type Chapter,
  type Character,
  type PlotEvent,
  type Scene,
  type ValidationIssue,
  type Version,
  type AuditLog,
  type WorkflowStep
} from '@/mock/data'

export const useProjectStore = defineStore('project', {
  state: () => ({
    projects: [] as Project[],
    currentProjectId: 0,
    chapters: [] as Chapter[],
    characters: [...mockCharacters] as Character[],
    plotEvents: [...mockPlotEvents] as PlotEvent[],
    scenes: [...mockScenes] as Scene[],
    yamlScript: mockYamlScript as string,
    validationIssues: [...mockValidationIssues] as ValidationIssue[],
    versions: [...mockVersions] as Version[],
    auditLogs: [] as AuditLog[],
    workflowSteps: [...mockWorkflowSteps] as WorkflowStep[],
    selectedSceneId: '' as string,
    selectedCharacterId: '' as string,
    loading: false
  }),
  getters: {
    currentProject: (state) => state.projects.find(p => p.id === state.currentProjectId),
    selectedScene: (state) => state.scenes.find(s => s.id === state.selectedSceneId),
    selectedCharacter: (state) => state.characters.find(c => c.id === state.selectedCharacterId),
    unresolvedIssues: (state) => state.validationIssues.filter(i => !i.resolved),
    highRiskIssues: (state) => state.validationIssues.filter(i => i.severity === 'high' && !i.resolved)
  },
  actions: {
    async fetchProjects() {
      try {
        const res: any = await projectApi.list({ page: 1, page_size: 50 })
        if (res.code === 0) {
          this.projects = res.data.projects || []
          if (this.projects.length > 0 && !this.currentProjectId) {
            this.currentProjectId = this.projects[0].id
          }
        }
      } catch (e) {
        console.error('获取项目列表失败', e)
      }
    },

    async fetchChapters(projectId: number) {
      try {
        const res: any = await chapterApi.list(projectId)
        if (res.code === 0) {
          this.chapters = res.data.chapters || []
        }
      } catch (e) {
        console.error('获取章节失败', e)
      }
    },

    async fetchScript(projectId: number) {
      try {
        const res: any = await scriptApi.get(projectId)
        if (res.code === 0) {
          this.yamlScript = res.data.yaml_content || ''
          try {
            const parsed: any = yaml.load(this.yamlScript)
            const script = parsed?.script
            if (Array.isArray(script?.characters)) {
              this.characters = script.characters
            }
            if (Array.isArray(script?.scenes)) {
              this.scenes = script.scenes.map((scene: any) => ({
                ...scene,
                risk_level: scene.risk_level || 'none'
              }))
            }
          } catch (e) {
            console.error('解析剧本失败', e)
          }
        }
      } catch (e) {
        console.error('获取剧本失败', e)
      }
    },

    async fetchCharacters(projectId: number) {
      const res: any = await resultApi.characters(projectId)
      if (res.code === 0) {
        this.characters = (res.data.characters || []).map((item: any) => ({
          id: item.character_key,
          name: item.name,
          aliases: item.aliases || [],
          role: item.role_type,
          description: item.description,
          personality: item.personality || [],
          relationships: item.relationships || [],
          source_trace: item.source_refs || [],
          confidence: item.confidence
        }))
      }
    },

    async fetchPlotEvents(projectId: number) {
      const res: any = await resultApi.plotEvents(projectId)
      if (res.code === 0) {
        this.plotEvents = (res.data.plot_events || []).map((item: any) => ({
          id: item.event_key,
          chapter_index: item.chapter_index,
          trigger: item.trigger_text,
          action: item.action_text,
          result: item.result_text,
          importance: item.importance,
          characters_involved: [],
          source_trace: item.source_refs || [],
          confidence: item.confidence
        }))
      }
    },

    async fetchVersions(projectId: number) {
      const res: any = await resultApi.versions(projectId)
      if (res.code === 0) {
        this.versions = res.data.versions || []
      }
    },

    async fetchValidationIssues(projectId: number) {
      const res: any = await resultApi.validationIssues(projectId)
      if (res.code === 0) {
        this.validationIssues = (res.data.issues || []).map((issue: any) => ({
          ...issue,
          resolved: Boolean(issue.resolved)
        }))
      }
    },

    async fetchAuditLogs(projectId: number) {
      try {
        const res: any = await auditApi.list(projectId, { page: 1, page_size: 50 })
        if (res.code === 0) {
          this.auditLogs = res.data.logs || []
        }
      } catch (e) {
        console.error('获取审计日志失败', e)
      }
    },

    async saveScript(projectId: number, yamlContent: string) {
      const res: any = await scriptApi.update(projectId, { yaml_content: yamlContent })
      if (res.code !== 0) {
        throw new Error(res.message || '保存剧本失败')
      }
      this.yamlScript = yamlContent
    },

    async validateScript(projectId: number) {
      const res: any = await scriptApi.validate(projectId)
      if (res.code !== 0) {
        throw new Error(res.message || '校验失败')
      }
      this.validationIssues = (res.data.issues || []).map((issue: any, index: number) => ({
        id: index + 1,
        ...issue,
        resolved: false
      }))
      return res.data
    },

    async createProject(title: string, description: string, mode: string) {
      const res: any = await projectApi.create({
        title,
        description,
        adaptation_mode: mode
      })
      if (res.code === 0) {
        await this.fetchProjects()
        return res.data.id
      }
      throw new Error(res.message || '创建项目失败')
    },

    async triggerGenerate(projectId: number) {
      const res: any = await taskApi.create(projectId, { task_type: 'full_generate' })
      if (res.code === 0) {
        return res.data.task_id
      }
      throw new Error(res.message || '创建任务失败')
    },

    async checkTaskStatus(taskId: number) {
      const res: any = await taskApi.status(taskId)
      if (res.code === 0) {
        return res.data
      }
      throw new Error('获取任务状态失败')
    },

    setCurrentProject(id: number) {
      this.currentProjectId = id
    },
    selectScene(id: string) {
      this.selectedSceneId = id
    },
    selectCharacter(id: string) {
      this.selectedCharacterId = id
    },
    resolveIssue(id: number) {
      const issue = this.validationIssues.find(i => i.id === id)
      if (issue) issue.resolved = true
    },
    updateYaml(content: string) {
      this.yamlScript = content
    }
  }
})
