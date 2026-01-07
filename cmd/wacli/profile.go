package main

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/steipete/wacli/internal/out"
)

func newProfileCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "profile",
		Short: "Profile management",
	}
	cmd.AddCommand(newProfilePictureCmd(flags))
	return cmd
}

func newProfilePictureCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "picture",
		Short: "Manage profile picture",
	}
	cmd.AddCommand(newProfilePictureSetCmd(flags))
	cmd.AddCommand(newProfilePictureRemoveCmd(flags))
	return cmd
}

func newProfilePictureSetCmd(flags *rootFlags) *cobra.Command {
	var imagePath string

	cmd := &cobra.Command{
		Use:   "set",
		Short: "Set your profile picture",
		RunE: func(cmd *cobra.Command, args []string) error {
			if imagePath == "" {
				return fmt.Errorf("--image is required")
			}

			ctx, cancel := withTimeout(context.Background(), flags)
			defer cancel()

			a, lk, err := newApp(ctx, flags, true, false)
			if err != nil {
				return err
			}
			defer closeApp(a, lk)

			if err := a.EnsureAuthed(); err != nil {
				return err
			}

			if err := a.Connect(ctx, false, nil); err != nil {
				return err
			}

			// Read the image file
			imageData, err := os.ReadFile(imagePath)
			if err != nil {
				return fmt.Errorf("failed to read image: %w", err)
			}

			// Set the profile picture
			pictureID, err := a.WA().SetProfilePicture(ctx, imageData)
			if err != nil {
				return fmt.Errorf("failed to set profile picture: %w", err)
			}

			resp := map[string]any{
				"success":    true,
				"picture_id": pictureID,
			}
			if flags.asJSON {
				return out.WriteJSON(os.Stdout, resp)
			}
			fmt.Fprintf(os.Stdout, "Profile picture set successfully (ID: %s)\n", pictureID)
			return nil
		},
	}

	cmd.Flags().StringVar(&imagePath, "image", "", "path to image file (JPEG recommended)")
	_ = cmd.MarkFlagRequired("image")
	return cmd
}

func newProfilePictureRemoveCmd(flags *rootFlags) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "remove",
		Short: "Remove your profile picture",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx, cancel := withTimeout(context.Background(), flags)
			defer cancel()

			a, lk, err := newApp(ctx, flags, true, false)
			if err != nil {
				return err
			}
			defer closeApp(a, lk)

			if err := a.EnsureAuthed(); err != nil {
				return err
			}

			if err := a.Connect(ctx, false, nil); err != nil {
				return err
			}

			// Pass nil to remove the profile picture
			result, err := a.WA().SetProfilePicture(ctx, nil)
			if err != nil {
				return fmt.Errorf("failed to remove profile picture: %w", err)
			}

			resp := map[string]any{
				"success": true,
				"result":  result,
			}
			if flags.asJSON {
				return out.WriteJSON(os.Stdout, resp)
			}
			fmt.Fprintf(os.Stdout, "Profile picture removed successfully\n")
			return nil
		},
	}

	return cmd
}
