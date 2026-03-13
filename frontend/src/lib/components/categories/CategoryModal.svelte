<script lang="ts">
    import {
        categoryService,
        type Category,
        type CategoryCreate,
        type CategoryType,
    } from "$lib/services/categories";
    import { toast } from "$lib/stores/toast";
    import Button from "$lib/components/ui/Button.svelte";

    let { show, category, categoryTypes, categories, onSave, onCancel } =
        $props<{
            show: boolean;
            category: Partial<Category>;
            categoryTypes: CategoryType[];
            categories: Category[];
            onSave: () => void;
            onCancel: () => void;
        }>();

    let isSaving = $state(false);
    let errors = $state<Record<string, string>>({});

    // Filter categories that can be parents (same type, not itself, and not a child of itself - though 2 levels simplified means just same type and not itself)
    let parentOptions = $derived(
        categories.filter(
            (c: Category) =>
                c.type_id === category.type_id &&
                c.id !== category.id &&
                !c.parent_id, // For 2 levels, parents cannot have a parent
        ),
    );

    const COLORS = [
        "#3B82F6",
        "#10B981",
        "#F59E0B",
        "#EF4444",
        "#8B5CF6",
        "#EC4899",
        "#64748B",
        "#06B6D4",
        "#F97316",
        "#0ea5e9",
    ];

    const ICONS = [
        "label",
        "folder",
        "flag",
        "star",
        "bookmark",
        "sell",
        "category",
        "layers",
        "tag",
    ];

    function validate() {
        errors = {};
        if (!category.name?.trim()) errors.name = "Category Name is required";
        if (!category.type_id) errors.type_id = "Please select a category type";
        return Object.keys(errors).length === 0;
    }

    async function handleSave() {
        if (!validate()) return;
        isSaving = true;
        try {
            if (category.id) {
                await categoryService.updateCategory(category.id, category);
                toast.success("Category updated successfully");
            } else {
                await categoryService.createCategory(
                    category as CategoryCreate,
                );
                toast.success("Category added successfully");
            }
            onSave();
        } catch (error: any) {
            toast.error(error.message || "Error saving category");
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
                    {category.id ? "Update Category" : "Add New Category"}
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
                        for="cat-type"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                    >
                        Category Type <span class="text-rose-500">*</span>
                    </label>
                    <select
                        id="cat-type"
                        bind:value={category.type_id}
                        class="w-full px-4 py-2.5 bg-slate-50 border {errors.type_id
                            ? 'border-rose-300'
                            : 'border-slate-200'} rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                    >
                        <option value={0}>-- Select category type --</option>
                        {#each categoryTypes as ct}
                            <option value={ct.id}>{ct.name}</option>
                        {/each}
                    </select>
                </div>

                <div>
                    <label
                        for="cat-name"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                    >
                        Category Name <span class="text-rose-500">*</span>
                        {#if errors.name}<span
                                class="text-rose-500 normal-case font-normal italic ml-2"
                                >{errors.name}</span
                            >{/if}
                    </label>
                    <input
                        id="cat-name"
                        bind:value={category.name}
                        type="text"
                        class="w-full px-4 py-2.5 bg-slate-50 border {errors.name
                            ? 'border-rose-300'
                            : 'border-slate-200'} rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                        placeholder="e.g. Important, Urgent..."
                    />
                </div>

                <div>
                    <label
                        for="cat-parent"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                    >
                        Parent Category
                    </label>
                    <select
                        id="cat-parent"
                        bind:value={category.parent_id}
                        class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                    >
                        <option value={undefined}>-- None (Top Level) --</option
                        >
                        {#each parentOptions as po}
                            <option value={po.id}>{po.name}</option>
                        {/each}
                    </select>
                </div>

                <div role="group" aria-labelledby="color-label">
                    <label
                        id="color-label"
                        class="block text-xs font-bold text-slate-500 uppercase mb-2 ml-1"
                        >Color</label
                    >
                    <div class="flex flex-wrap gap-2 px-1">
                        {#each COLORS as color}
                            <button
                                class="w-8 h-8 rounded-full border-2 {category.color ===
                                color
                                    ? 'border-primary ring-2 ring-primary/20'
                                    : 'border-transparent'}"
                                style="background-color: {color}"
                                onclick={() => (category.color = color)}
                                type="button"
                                aria-label="Select color {color}"
                            ></button>
                        {/each}
                    </div>
                </div>

                <div role="group" aria-labelledby="icon-label">
                    <label
                        id="icon-label"
                        class="block text-xs font-bold text-slate-500 uppercase mb-2 ml-1"
                        >Icon</label
                    >
                    <div class="flex flex-wrap gap-2 px-1">
                        {#each ICONS as icon}
                            <button
                                class="w-10 h-10 flex items-center justify-center rounded-xl bg-slate-50 border-2 {category.icon ===
                                icon
                                    ? 'border-primary text-primary'
                                    : 'border-transparent text-slate-400'} hover:bg-slate-100 transition-all font-display"
                                onclick={() => (category.icon = icon)}
                                type="button"
                                aria-label="Select icon {icon}"
                            >
                                <span class="material-symbols-outlined"
                                    >{icon}</span
                                >
                            </button>
                        {/each}
                    </div>
                </div>

                <div>
                    <label
                        for="cat-desc"
                        class="block text-xs font-bold text-slate-500 uppercase mb-1.5 ml-1"
                        >Description</label
                    >
                    <textarea
                        id="cat-desc"
                        bind:value={category.description}
                        rows="2"
                        class="w-full px-4 py-2.5 bg-slate-50 border border-slate-200 rounded-xl focus:ring-2 focus:ring-primary/20 outline-none transition-all"
                        placeholder="Enter detailed category description..."
                    ></textarea>
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
