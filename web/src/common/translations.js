import { int } from './'

export default {
  /* APP */
  app_drag_drop_message:                int._t('Drop you files here'),
  app_drag_drop_success:                int._t('Uploaded'),
  app_drag_drop_error:                  int._t('An error occured, please try again'),

  /* DOWNLOADS */
  downloads_empty_state_title:          int._t('No downloads whatsover'),
  downloads_empty_state_message:        int._t('You can setup your sync app by clicking the button bellow'),
  downloads_empty_state_cta_label:      int._t('Settings'),

  /* SETTINGS */
  settings_poll_interval_message:       int._t('Poll Interval'),
  settings_source_folder_label:         int._t('Source folder'),
  settings_source_folder_action_label:  int._t('Change'),
  settings_dest_folder_label:           int._t('Destination folder'),
  settings_dest_folder_action_label:    int._t('Change'),
  settings_simultaneous_download_label: int._t('Simultaneous download'),
  settings_segments_perfile_label:      int._t('Segments per file'),
  settings_save_success_message:        int._t('Saved'),
  settings_delete_source_message:       int._t('Delete original file'),
  settings_delete_source_hint:          int._t('After a successful download'),
  settings_save_label:                  int._t('Save'),

  /* FILE TREE */
  filetree_title:                       int._t('Please choose a folder'),
  filetree_action_cancel_label:         int._t('Cancel'),

  /* NEW TRANSFER */
  new_transfer_no_link_error:           int._t('There is no link to download'),
  new_transfer_invalid_link_error:      int._t('Couldn\'t find anything to download there. We accept torrent hashes, magnet links, links to torrents and some video page links.'),
  new_transfer_success_message:         int._t('Your transfers have been started.'),
}
