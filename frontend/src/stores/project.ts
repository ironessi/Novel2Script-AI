import { defineStore } from 'pinia'
import {
  mockProjects,
  mockChapters,
  mockCharacters,
  mockPlotEvents,
  mockScenes,
  mockYamlScript,
  mockValidationIssues,
  mockVersions,
  mockAuditLogs,
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
    projects: [...mockProjects] as Project[],
    currentProjectId: 1,
    chapters: [...mockChapters] as Chapter[],
    characters: [...mockCharacters] as Character[],
    plotEvents: [...mockPlotEvents] as PlotEvent[],
    scenes: [...mockScenes] as Scene[],
    yamlScript: mockYamlScript as string,
    validationIssues: [...mockValidationIssues] as ValidationIssue[],
    versions: [...mockVersions] as Version[],
    auditLogs: [...mockAuditLogs] as AuditLog[],
    workflowSteps: [...mockWorkflowSteps] as WorkflowStep[],
    selectedSceneId: '' as string,
    selectedCharacterId: '' as string
  }),
  getters: {
    currentProject: (state) => state.projects.find(p => p.id === state.currentProjectId),
    selectedScene: (state) => state.scenes.find(s => s.id === state.selectedSceneId),
    selectedCharacter: (state) => state.characters.find(c => c.id === state.selectedCharacterId),
    unresolvedIssues: (state) => state.validationIssues.filter(i => !i.resolved),
    highRiskIssues: (state) => state.validationIssues.filter(i => i.severity === 'high' && !i.resolved)
  },
  actions: {
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
