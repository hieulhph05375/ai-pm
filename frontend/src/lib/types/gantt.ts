export interface GanttItem {
    id: string | number;
    title: string;
    startDate: string | null;
    dueDate: string | null;
    progress: number;
    type?: 'phase' | 'milestone' | 'task' | 'sub-task';
    level?: number;
    parent_id?: string | number | null;
}

export interface GanttDependency {
    id: string | number;
    fromId: string | number;
    toId: string | number;
    type: 'FS' | 'SS' | 'FF' | 'SF';
}
