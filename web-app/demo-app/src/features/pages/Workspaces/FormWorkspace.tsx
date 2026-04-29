import { Button } from '@/components/ui/button'
import { Input } from '@/components/ui/input'
import {
    Select,
    SelectContent,
    SelectItem,
    SelectTrigger,
    SelectValue,
} from '@/components/ui/select'
import { cn } from '@/lib/utils'

import { listTypeTimeDuration, useCreateWorkspace } from '@/features/hooks/useCreateWorkspace'
import { Textarea } from '@/components/ui/textarea'

const listTemplate = [
    { id: 'id-01', name: 'VsCode', icon: '📝' },
    { id: 'id-02', name: 'Postgres', icon: '🐘' },
    { id: 'id-03', name: 'Blank', icon: '📦' },
]

export function FormWorkspaces() {
    const { register, handleSubmit, setValue, watch } = useCreateWorkspace()

    const selectedUnit = watch('typeTimeDuration')

    return (
        <form onSubmit={handleSubmit} className="w-full ">
            <div className="rounded-xl border border-border bg-card p-6 space-y-5">

                {/* Header */}
                <h2 className="text-lg uppercase tracking-widest text-muted-foreground font-medium">
                    New workspace
                </h2>

                {/* Workspace Name */}
                <div className="space-y-3 flex flex-col">
                    <label className="text-sm text-muted-foreground">
                        Workspace name
                    </label>
                    <Input
                        {...register('name')}
                        placeholder="e.g. my-dev-env"
                        className="border-2 border-blue-300/20"
                    />
                </div>

                {/* Description */}
                <div className="space-y-3 flex flex-col">
                    <label className="text-sm text-muted-foreground">
                        Description
                    </label>
                    <Textarea
                        {...register('description')}
                        placeholder="Optional description"
                        className="border-2 border-blue-300/20"
                    />
                </div>

                {/* Template selector */}
                <div className="space-y-3 flex flex-col">
                    <label className="text-sm text-muted-foreground">Template</label>
                    <input type="hidden" {...register('templateId')} />
                    <Select
                        onValueChange={(val) =>
                            setValue('templateId', val, { shouldValidate: true })
                        }
                        defaultValue=""
                    >
                        <SelectTrigger className="w-full border-2 border-blue-300/20">
                            <SelectValue placeholder="Choose a template..." />
                        </SelectTrigger>
                        <SelectContent>
                            {listTemplate.map((item) => (
                                <SelectItem key={item.id} value={item.id}>
                                    {item.name}
                                </SelectItem>
                            ))}
                        </SelectContent>
                    </Select>
                </div>

                {/* Duration */}
                <div className="grid grid-cols-2 gap-3">
                    <div className="space-y-3 flex flex-col">
                        <label className="text-sm text-muted-foreground">
                            Duration
                        </label>
                        <Input
                            {...register('timeDuration')}
                            type="number"
                            min={1}
                            placeholder="e.g. 8"
                            className="border-2 border-blue-300/20"
                        />
                    </div>

                    <div className="space-y-3 flex flex-col">
                        <label className="text-sm text-muted-foreground">
                            Unit
                        </label>
                        {/* hidden input so react-hook-form picks up the value */}
                        <input type="hidden" {...register('typeTimeDuration')} />
                        <div className="flex gap-1.5">
                            {listTypeTimeDuration.map((unit) => (
                                <button
                                    key={unit}
                                    type="button"
                                    onClick={() => setValue('typeTimeDuration', unit, { shouldValidate: true })}
                                    className={cn(
                                        'flex-1 rounded-md border px-2 py-2 text-xs transition-all',
                                        'hover:border-foreground/40 focus-visible:outline-none focus-visible:ring-2 focus-visible:ring-ring',
                                        selectedUnit === unit
                                            ? 'border-foreground bg-blue-300/20 font-medium text-foreground'
                                            : 'border-border bg-muted/40 text-muted-foreground'
                                    )}
                                >
                                    {unit}
                                </button>
                            ))}
                        </div>
                    </div>
                    {/* Workspace limit CPU */}
                    <div className="space-y-3 flex flex-col">
                        <label className="text-sm text-muted-foreground">
                            Request CPU
                        </label>
                        <Input
                            type='number'
                            {...register('reqCpu')}
                            placeholder="0.30"
                            className="border-2 border-blue-300/20"
                        />
                    </div>
                      <div className="space-y-3 flex flex-col">
                        <label className="text-sm text-muted-foreground">
                            Request RAM
                        </label>
                        <Input
                            type='number'
                            {...register('reqRam')}
                            placeholder="0.30"
                            className="border-2 border-blue-300/20"
                        />
                    </div>
                      <div className="space-y-3 flex flex-col">
                        <label className="text-sm text-muted-foreground">
                            Limit CPU
                        </label>
                        <Input
                            type='number'
                            {...register('limitCpu')}
                            placeholder="0.30"
                            className="border-2 border-blue-300/20"
                        />
                    </div>
                      <div className="space-y-3 flex flex-col">
                        <label className="text-sm text-muted-foreground">
                            Limit RAM
                        </label>
                        <Input
                            type='number'
                            {...register('limitRam')}
                            placeholder="0.30"
                            className="border-2 border-blue-300/20"
                        />
                    </div>
                </div>

                {/* Submit */}
                <div className="flex pt-1">
                    <Button type="submit" className=" w-full px-6">
                        Create workspace
                    </Button>
                </div>

            </div>
        </form>
    )
}