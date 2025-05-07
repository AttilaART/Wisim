export type window = {
  id: number,
  name: string,
  z_index: number,
  hidden: boolean,
  close_window: () => void
}
export const windows: window[] = $state([])

export function new_window(name: string, close: () => void): number {
  windows.push({ id: windows.length, name, z_index: windows.length + 1, hidden: false, close_window: close })
  return windows[windows.length - 1].id
}
export function move_window_to_top(window_id: number) {
  let current_z_index = windows[window_id].z_index
  for (let i in windows) {
    if (windows[i].z_index > current_z_index) {
      windows[i].z_index -= 1
    }
  }
  windows[window_id].z_index = windows.length
}

export function delete_window(window_id: number) {
  let deleted_window_z_index: number = windows[window_id].z_index
  windows[window_id].close_window()

  windows.splice(window_id, 1);

  for (let i in windows) {
    if (windows[i].z_index > deleted_window_z_index) {
      windows[i].z_index -= 1
    }
  }
}
