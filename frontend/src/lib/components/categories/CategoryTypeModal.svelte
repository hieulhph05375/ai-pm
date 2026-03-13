<script lang="ts">
    import {
        categoryService,
        type CategoryType,
        type CategoryTypeCreate,
    } from "$lib/services/categories";
    import { toast } from "$lib/stores/toast";
    import Button from "$lib/components/ui/Button.svelte";

    let { show, categoryType, onSave, onCancel } = $props<{
        show: boolean;
        categoryType: Partial<CategoryType>;
        onSave: () => void;
        onCancel: () => void;
    }>();

    let isSaving = $state(false);
    let errors = $state<Record<string, string>>({});

    function validate() {
        errors = {};
        if (!categoryType.name?.trim()) errors.name = "Type Name is required";
        if (!categoryType.code?.trim()) errors.code = "Type Code is required";
        return Object.keys(errors).length === 0;
    }

    async function handleSave() {
        if (!validate()) return;
        isSaving = true;
        try {
            if (categoryType.id) {
                await categoryService.updateType(categoryType.id, categoryType);
                toast.success("Category type updated successfully");
            } else {
                await categoryService.createType(
                    categoryType as CategoryTypeCreate,
                );
                toast.success("Category type added successfully");
            }
            onSave();
        } catch (error: any) {
            toast.error(error.message || "Error saving category type");
        } finally {
            isSaving = false;
        }
    }
</script>

{#if show}
    <div class="fixed inset-0 z-[110] flex items-center justify-center p-4">
        <button
            class="absolute inset-0 bg-slate-900/40 backdrop-blur-sm border-none w-full h-full cursor-default"
            onclick={onCancel}
            aria-label="Close modal"
        ></button>
        <div
            class="bg-white rounded-2xl w-full max-w-md relative shadow-2xl animate-in fade-in zoom-in duration-200 overflow-hidden"
        >
            <div
                class="px-6 py-4 border-b border-slate-100 flex justify-between items-center bg-slate-50/50"
            >
                <h3 class="font-display font-bold text-lg text-slate-900">
                    {categoryType.id
                        ? "Update Category Type"
                        : "Add New Category Type"}
                </h3>
                <button
                    onclick={onCancel}
                    class="text-slate-400 hover:text-slate-600"
                >
                    <span class="material-symbols-outlined">close</span>
                </button>
            </div>

            <div class="p-6 space-y-4">
                <div>
                    <label
                        for="type-name"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                    >
                        Type Name <span class="text-rose-500">*</span>
                        {#if errors.name}<span
                                class="text-rose-500 normal-case font-normal italic ml-2"
                                >{errors.name}</span
                            >{/if}
                    </label>
                    <input
                        id="type-name"
                        bind:value={categoryType.name}
                        type="text"
                        class="w-full px-4 py-2.5 bg-slate-50 border {errors.name
                            ? 'border-rose-300'
                            : 'border-slate-200'} rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                        placeholder="e.g. Priority, Department..."
                    />
                </div>

                <div>
                    <label
                        for="type-code"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                    >
                        Code <span class="text-rose-500">*</span>
                        {#if errors.code}<span
                                class="text-rose-500 normal-case font-normal italic ml-2"
                                >{errors.code}</span
                            >{/if}
                    </label>
                    <input
                        id="type-code"
                        bind:value={categoryType.code}
                        type="text"
                        class="w-full px-4 py-2.5 bg-slate-50 border {errors.code
                            ? 'border-rose-300'
                            : 'border-slate-200'} rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                        placeholder="e.g. PRIORITY, DEPT..."
                        disabled={!!categoryType.id}
                    />
                </div>

                <div>
                    <label
                        for="type-desc"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                        >Description</label
                    >
                    <textarea
                        id="type-desc"
                        bind:value={categoryType.description}
                        rows="3"
                        class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                        placeholder="Brief description of this category type..."
                    ></textarea>
                </div>

                <div class="flex items-center gap-2 px-1">
                    <input
                        type="checkbox"
                        id="type-active"
                        bind:checked={categoryType.is_active}
                        class="w-4 h-4 rounded border-slate-300 text-primary focus:ring-primary"
                    />
                    <label
                        for="type-active"
                        class="text-sm font-medium text-slate-700 font-display"
                        >Active</label
                    >
                </div>
            </div>

            <div
                class="p-6 bg-slate-50/50 border-t border-slate-100 flex justify-end gap-3"
            >
                <button
                    onclick={onCancel}
                    class="px-5 py-2.5 rounded-xl font-medium text-slate-600 hover:bg-slate-200 transition-all"
                    >Cancel</button
                >
                <Button loading={isSaving} onclick={handleSave}>Save</Button>
            </div>
        </div>
    </div>
{/if}
