export type Window = {
  id: number,
  name: string,
  z_index: number,
  hidden: boolean,
  onClose: () => void
}
export const windows: Window[] = $state([])

export function get_window_by_id(window_id: number): { index: number, window: Window } {
  for (let i in windows) {
    if (windows[i].id == window_id) {
      return { index: Number(i), window: windows[i] }
    }
  }

  throw new Error(`no window with id==${window_id} found.`)
}

export function window_exists(window_id: number): boolean {
  try {
    get_window_by_id(window_id)
    return true
  } catch {
    return false
  }
}

export function new_window(name: string, close: () => void): number {
  let id: number = 0
  for (let i = 0; i <= windows.length; i++) {
    let index_taken: boolean = false
    for (let k in windows) {
      if (index_taken) break
      if (Number(i) == windows[k].id) {
        index_taken = true
      }
    }

    if (!index_taken) {
      id = Number(i)
      break
    }
  }
  windows.push({ id: id, name, z_index: windows.length + 1, hidden: false, onClose: close })
  return id
}
export function move_window_to_top(window_id: number) {
  let current_z_index = windows[get_window_by_id(window_id).index].z_index
  let current_index = get_window_by_id(window_id).index
  for (let i in windows) {
    if (windows[i].z_index >= current_z_index && Number(i) != current_index) {
      windows[i].z_index -= 1
    }
  }
  windows[current_index].z_index = windows.length
  console.log($state.snapshot(windows))
}

export function delete_window(window_id: number) {
  let deleted_window_z_index: number = windows[get_window_by_id(window_id).index].z_index
  windows[get_window_by_id(window_id).index].onClose()

  windows.splice(get_window_by_id(window_id).index, 1);

  for (let i in windows.keys()) {
    if (windows[i].z_index > deleted_window_z_index) {
      windows[i].z_index -= 1
    }
  }
}
